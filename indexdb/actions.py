import hashlib
from itertools import chain
import logging
import json
from typing import List, Dict
import uuid


from bugout.data import BugoutSearchResult


from .models import AbiJobs
from .settings import (
    MOONSTREAM_ADMIN_ACCESS_TOKEN,
    MOONSTREAM_MOONWORM_TASKS_JOURNAL_ID,
)

from .settings import bugout_client as bc
from .db import yield_db_session_ctx


logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)


def get_all_entries_from_search(
    journal_id: str, search_query: str, limit: int, token: str, content: bool = True
) -> List[BugoutSearchResult]:
    """
    Get all required entries from journal using search interface
    """
    offset = 0

    results: List[BugoutSearchResult] = []

    existing_methods = bc.search(
        token=token,
        journal_id=journal_id,
        query=search_query,
        content=content,
        timeout=10.0,
        limit=limit,
        offset=offset,
    )
    results.extend(existing_methods.results)  # type: ignore

    if len(results) != existing_methods.total_results:
        for offset in range(limit, existing_methods.total_results, limit):
            existing_methods = bc.search(
                token=token,
                journal_id=journal_id,
                query=search_query,
                content=content,
                timeout=10.0,
                limit=limit,
                offset=offset,
            )
            results.extend(existing_methods.results)  # type: ignore

    return results


def migrate_moonworm_tasks(
    subscription_type: str,
    entries_limit: int = 100,
) -> None:
    """
    Get list of subscriptions loads abi and apply them as moonworm tasks if it not exist
    """

    # get all existing moonworm tasks address and list of selectors
    parsed_subscriptions: Dict[str, List[str]] = dict()

    with yield_db_session_ctx() as session:

        existing_selectors = session.query(AbiJobs.address, AbiJobs.abi_selector).all()

        for address, abi_selector in existing_selectors:
            if address not in parsed_subscriptions:
                parsed_subscriptions[address] = []
            parsed_subscriptions[address].append(abi_selector)

    abi_jobs = []

    user_id = uuid.UUID("00000000-0000-0000-0000-000000000000")
    customer_id = uuid.UUID("00000000-0000-0000-0000-000000000000")

    abi_selector = None

    print(f"subscription_type: {subscription_type}")

    test = []

    try:
        entries = get_all_entries_from_search(
            journal_id=MOONSTREAM_MOONWORM_TASKS_JOURNAL_ID,
            search_query=f"tag:subscription_type:{subscription_type}",
            limit=entries_limit,  # load per request
            token=MOONSTREAM_ADMIN_ACCESS_TOKEN,
        )

        for entry in entries:

            ## get abi
            try:
                abi = json.loads(entry.content)
            except Exception as e:
                logger.warn(
                    f"Unable to parse abi from task: {entry.entry_url.split()[-1]}: {e}"
                )
                continue

            for tag in entry.tags:
                if tag.startswith("abi_selector:"):
                    abi_selector = tag.split(":")[1]
                if tag.startswith("address:"):
                    address = tag.split(":")[1]
                if tag.startswith("abi_name:"):
                    abi_name = tag.split(":")[1]
                if tag.startswith("historical_crawl_status:"):
                    historical_crawl_status = tag.split(":")[1]
                if tag.startswith("status:"):
                    status = tag.split(":")[1]
                if tag.startswith("progress:"):
                    try:
                        progress = int(float(tag.split(":")[1]))
                    except:
                        logger.warn(
                            f"Unable to parse progress from task: {entry.entry_url.split()[-1]}"
                        )
                        raise
                if tag.startswith("moonworm_task_pickedup:"):
                    moonworm_task_pickedup = tag.split(":")[1] == "True"

            if address not in parsed_subscriptions:
                parsed_subscriptions[address] = []

            if abi_selector is None:
                ### generate abi_selector

                hash = hashlib.md5(json.dumps(abi).encode("utf-8")).hexdigest()
                abi_selector = f"0x{hash[:8]}"
                logger.info(
                    f"Generated abi_selector: {abi_selector} for address {address}"
                )

            if abi_selector not in parsed_subscriptions[address]:
                parsed_subscriptions[address].append(abi_selector)
            else:
                logger.warning(
                    f"Duplicate selector {abi_selector} for address {address}"
                )
                continue

            if address == "0x36f2119F866B27032cDf3A8CDf19FB3b63BaD281":

                test.append(
                    {
                        "address": address,
                        "user_id": str(user_id),
                        "customer_id": str(customer_id),
                        "abi_selector": abi_selector,
                        "chain": subscription_type,
                        "abi_name": abi_name,
                    }
                )

            abi_jobs.append(
                AbiJobs(
                    address=address,
                    user_id=user_id,
                    customer_id=customer_id,
                    abi_selector=abi_selector,
                    chain=subscription_type,
                    abi_name=abi_name,
                    status=status,
                    historical_crawl_status=historical_crawl_status,
                    progress=progress,
                    moonworm_task_pickedup=moonworm_task_pickedup,
                    abi=json.dumps(abi),
                )
            )

    except Exception as e:
        logger.error(f"Error parse moonworm tasks: {str(e)}")

    with open("test.json", "w") as f:
        f.write(json.dumps(test, indent=4))

    with yield_db_session_ctx() as session:
        session.add_all(abi_jobs)
        session.commit()
        print(f"moonworm_jobs {session.query(AbiJobs).count()}")

    # get all subscriptions
