"""tx call and event unique indexes

Revision ID: d2ceff33be47
Revises: e9f640a2b45b
Create Date: 2024-05-17 16:35:56.059519

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = 'd2ceff33be47'
down_revision: Union[str, None] = 'e9f640a2b45b'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    # uq_amoy_labels
    op.create_index('uk_amoy_labels_tx_hash_log_idx_evt', 'amoy_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_amoy_labels_tx_hash_tx_call', 'amoy_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_arbitrum_nova_labels
    op.create_index('uk_arbitrum_nova_labels_tx_hash_log_idx_evt', 'arbitrum_nova_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_arbitrum_nova_labels_tx_hash_tx_call', 'arbitrum_nova_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_arbitrum_one_labels
    op.create_index('uk_arbitrum_one_labels_tx_hash_log_idx_evt', 'arbitrum_one_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_arbitrum_one_labels_tx_hash_tx_call', 'arbitrum_one_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_arbitrum_sepolia_labels
    op.create_index('uk_arbitrum_sepolia_labels_tx_hash_log_idx_evt', 'arbitrum_sepolia_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_arbitrum_sepolia_labels_tx_hash_tx_call', 'arbitrum_sepolia_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_avalanche_fuji_labels
    op.create_index('uk_avalanche_fuji_labels_tx_hash_log_idx_evt', 'avalanche_fuji_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_avalanche_fuji_labels_tx_hash_tx_call', 'avalanche_fuji_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_avalanche_labels
    op.create_index('uk_avalanche_labels_tx_hash_log_idx_evt', 'avalanche_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_avalanche_labels_tx_hash_tx_call', 'avalanche_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_base_labels
    op.create_index('uk_base_labels_tx_hash_log_idx_evt', 'base_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_base_labels_tx_hash_tx_call', 'base_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_blast_labels
    op.create_index('uk_blast_labels_tx_hash_log_idx_evt', 'blast_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_blast_labels_tx_hash_tx_call', 'blast_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_blast_sepolia_labels
    op.create_index('uk_blast_sepolia_labels_tx_hash_log_idx_evt', 'blast_sepolia_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_blast_sepolia_labels_tx_hash_tx_call', 'blast_sepolia_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_ethereum_labels
    op.create_index('uk_ethereum_labels_tx_hash_log_idx_evt', 'ethereum_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_ethereum_labels_tx_hash_tx_call', 'ethereum_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_mumbai_labels
    op.create_index('uk_mumbai_labels_tx_hash_log_idx_evt', 'mumbai_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_mumbai_labels_tx_hash_tx_call', 'mumbai_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_polygon_labels
    op.create_index('uk_polygon_labels_tx_hash_log_idx_evt', 'polygon_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_polygon_labels_tx_hash_tx_call', 'polygon_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_proofofplay_apex_labels
    op.create_index('uk_proofofplay_apex_labels_tx_hash_log_idx_evt', 'proofofplay_apex_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_proofofplay_apex_labels_tx_hash_tx_call', 'proofofplay_apex_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_sepolia_labels
    op.create_index('uk_sepolia_labels_tx_hash_log_idx_evt', 'sepolia_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_sepolia_labels_tx_hash_tx_call', 'sepolia_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_starknet_labels
    op.create_index('uk_starknet_labels_tx_hash_log_idx_evt', 'starknet_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_starknet_labels_tx_hash_tx_call', 'starknet_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_starknet_sepolia_labels
    op.create_index('uk_starknet_sepolia_labels_tx_hash_log_idx_evt', 'starknet_sepolia_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_starknet_sepolia_labels_tx_hash_tx_call', 'starknet_sepolia_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_xai_labels
    op.create_index('uk_xai_labels_tx_hash_log_idx_evt', 'xai_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_xai_labels_tx_hash_tx_call', 'xai_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # uq_xai_sepolia_labels
    op.create_index('uk_xai_sepolia_labels_tx_hash_log_idx_evt', 'xai_sepolia_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_xai_sepolia_labels_tx_hash_tx_call', 'xai_sepolia_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # xdai_labels
    op.create_index('uk_xdai_labels_tx_hash_log_idx_evt', 'xdai_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_xdai_labels_tx_hash_tx_call', 'xdai_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # zksync_era_labels
    op.create_index('uk_zksync_era_labels_tx_hash_log_idx_evt', 'zksync_era_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_zksync_era_labels_tx_hash_tx_call', 'zksync_era_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # zksync_era_sepolia_labels
    op.create_index('uk_zksync_era_sepolia_labels_tx_hash_log_idx_evt', 'zksync_era_sepolia_labels', ['transaction_hash', 'log_index'], unique=True, postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.create_index('uk_zksync_era_sepolia_labels_tx_hash_tx_call', 'zksync_era_sepolia_labels', ['transaction_hash'], unique=True, postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_index('uk_zksync_era_sepolia_labels_tx_hash_tx_call', table_name='zksync_era_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_zksync_era_sepolia_labels_tx_hash_log_idx_evt', table_name='zksync_era_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_zksync_era_labels_tx_hash_tx_call', table_name='zksync_era_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_zksync_era_labels_tx_hash_log_idx_evt', table_name='zksync_era_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_xdai_labels_tx_hash_tx_call', table_name='xdai_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_xdai_labels_tx_hash_log_idx_evt', table_name='xdai_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_xai_sepolia_labels_tx_hash_tx_call', table_name='xai_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_xai_sepolia_labels_tx_hash_log_idx_evt', table_name='xai_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_xai_labels_tx_hash_tx_call', table_name='xai_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_xai_labels_tx_hash_log_idx_evt', table_name='xai_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_starknet_sepolia_labels_tx_hash_tx_call', table_name='starknet_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_starknet_sepolia_labels_tx_hash_log_idx_evt', table_name='starknet_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_starknet_labels_tx_hash_tx_call', table_name='starknet_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_starknet_labels_tx_hash_log_idx_evt', table_name='starknet_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_sepolia_labels_tx_hash_tx_call', table_name='sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_sepolia_labels_tx_hash_log_idx_evt', table_name='sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_proofofplay_apex_labels_tx_hash_tx_call', table_name='proofofplay_apex_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_proofofplay_apex_labels_tx_hash_log_idx_evt', table_name='proofofplay_apex_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_polygon_labels_tx_hash_tx_call', table_name='polygon_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_polygon_labels_tx_hash_log_idx_evt', table_name='polygon_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_mumbai_labels_tx_hash_tx_call', table_name='mumbai_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_mumbai_labels_tx_hash_log_idx_evt', table_name='mumbai_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_ethereum_labels_tx_hash_tx_call', table_name='ethereum_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_ethereum_labels_tx_hash_log_idx_evt', table_name='ethereum_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_blast_sepolia_labels_tx_hash_tx_call', table_name='blast_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_blast_sepolia_labels_tx_hash_log_idx_evt', table_name='blast_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_blast_labels_tx_hash_tx_call', table_name='blast_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_blast_labels_tx_hash_log_idx_evt', table_name='blast_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_base_labels_tx_hash_tx_call', table_name='base_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_base_labels_tx_hash_log_idx_evt', table_name='base_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_avalanche_labels_tx_hash_tx_call', table_name='avalanche_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_avalanche_labels_tx_hash_log_idx_evt', table_name='avalanche_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_avalanche_fuji_labels_tx_hash_tx_call', table_name='avalanche_fuji_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_avalanche_fuji_labels_tx_hash_log_idx_evt', table_name='avalanche_fuji_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_arbitrum_sepolia_labels_tx_hash_tx_call', table_name='arbitrum_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_arbitrum_sepolia_labels_tx_hash_log_idx_evt', table_name='arbitrum_sepolia_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_arbitrum_one_labels_tx_hash_tx_call', table_name='arbitrum_one_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_arbitrum_one_labels_tx_hash_log_idx_evt', table_name='arbitrum_one_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_arbitrum_nova_labels_tx_hash_tx_call', table_name='arbitrum_nova_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_arbitrum_nova_labels_tx_hash_log_idx_evt', table_name='arbitrum_nova_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    op.drop_index('uk_amoy_labels_tx_hash_tx_call', table_name='amoy_labels', postgresql_where=sa.text("label='seer' and label_type='tx_call'"))
    op.drop_index('uk_amoy_labels_tx_hash_log_idx_evt', table_name='amoy_labels', postgresql_where=sa.text("label='seer' and label_type='event'"))
    # ### end Alembic commands ###