import argparse

from .actions import migrate_moonworm_tasks

chain_to_sub_type = {
    "ethereum": "ethereum_smartcontract",
    "polygon": "polygon_smartcontract",
    "mumbai": "mumbai_smartcontract",
}


def migrate_moonworm_tasks_handler(args: argparse.Namespace) -> None:
    """
    Migrate moonworm tasks to moonstream.
    """
    ### print all args
    print(args)

    # Get the moonstream access token
    migrate_moonworm_tasks(
        subscription_type=chain_to_sub_type[args.blockchain],
        entries_limit=args.entries_limit,
    )


def main():
    parser = argparse.ArgumentParser(description="Migrate moonworm tasks to moonstream")
    parser.add_argument(
        "-b",
        "--blockchain",
        type=str,
        required=True,
        choices=["ethereum", "polygon", "mumbai"],
    )
    parser.add_argument(
        "--entries_limit",
        type=int,
        default=100,
        help="Limit of entries to migrate",
    )
    args = parser.parse_args()
    migrate_moonworm_tasks_handler(args)


if __name__ == "__main__":
    main()
