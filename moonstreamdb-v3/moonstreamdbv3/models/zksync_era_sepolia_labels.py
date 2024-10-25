# zksync_era_sepolia_labels.py

from sqlalchemy import Index, text
from models.abstract_labels import EvmBasedLabel


class ZksyncEraSepoliaLabel(EvmBasedLabel):
    __tablename__ = "zksync_era_sepolia_labels"

    __table_args__ = (
        Index(
            "ix_zksync_era_sepolia_labels_addr_block_num",
            "address",
            "block_number",
            unique=False,
        ),
        Index(
            "ix_zksync_era_sepolia_labels_addr_block_ts",
            "address",
            "block_timestamp",
            unique=False,
        ),
        Index(
            "uk_zksync_era_sepolia_labels_tx_hash_tx_call",
            "transaction_hash",
            unique=True,
            postgresql_where=text("label='seer' and label_type='tx_call'"),
        ),
        Index(
            "uk_zksync_era_sepolia_labels_tx_hash_log_idx_evt",
            "transaction_hash",
            "log_index",
            unique=True,
            postgresql_where=text("label='seer' and label_type='event'"),
        ),
        Index(
            "uk_zksync_era_sepolia_labels_tx_hash_tx_call_raw",
            "transaction_hash",
            unique=True,
            postgresql_where=text("label='seer-raw' and label_type='tx_call'"),
        ),
        Index(
            "uk_zksync_era_sepolia_labels_tx_hash_log_idx_evt_raw",
            "transaction_hash",
            "log_index",
            unique=True,
            postgresql_where=text("label='seer-raw' and label_type='event'"),
        ),
    )
