"""Add game7 testnet

Revision ID: b89a8affe2c3
Revises: f19652e59bc5
Create Date: 2024-08-05 16:29:33.091636

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "b89a8affe2c3"
down_revision: Union[str, None] = "f19652e59bc5"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.create_table(
        "game7_testnet_blocks",
        sa.Column("l1_block_number", sa.BigInteger(), nullable=False),
        sa.Column("block_number", sa.BigInteger(), nullable=False),
        sa.Column("block_hash", sa.VARCHAR(length=256), nullable=False),
        sa.Column("block_timestamp", sa.BigInteger(), nullable=False),
        sa.Column("parent_hash", sa.VARCHAR(length=256), nullable=False),
        sa.Column("row_id", sa.BigInteger(), nullable=False),
        sa.Column("path", sa.Text(), nullable=False),
        sa.Column(
            "indexed_at",
            sa.DateTime(timezone=True),
            server_default=sa.text("TIMEZONE('utc', statement_timestamp())"),
            nullable=False,
        ),
        sa.PrimaryKeyConstraint("block_number", name=op.f("pk_game7_testnet_blocks")),
    )
    op.create_index(
        op.f("ix_game7_testnet_blocks_block_number"),
        "game7_testnet_blocks",
        ["block_number"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_blocks_block_timestamp"),
        "game7_testnet_blocks",
        ["block_timestamp"],
        unique=False,
    )
    op.create_table(
        "game7_testnet_reorgs",
        sa.Column("id", sa.UUID(), nullable=False),
        sa.Column("block_number", sa.BigInteger(), nullable=False),
        sa.Column("block_hash", sa.VARCHAR(length=256), nullable=False),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_game7_testnet_reorgs")),
    )
    op.create_index(
        op.f("ix_game7_testnet_reorgs_block_hash"),
        "game7_testnet_reorgs",
        ["block_hash"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_reorgs_block_number"),
        "game7_testnet_reorgs",
        ["block_number"],
        unique=False,
    )
    op.create_table(
        "game7_testnet_transactions",
        sa.Column("block_number", sa.BigInteger(), nullable=False),
        sa.Column("hash", sa.VARCHAR(length=256), nullable=False),
        sa.Column("from_address", sa.LargeBinary(length=20), nullable=False),
        sa.Column("to_address", sa.LargeBinary(length=20), nullable=False),
        sa.Column("selector", sa.VARCHAR(length=256), nullable=True),
        sa.Column("type", sa.Integer(), nullable=True),
        sa.Column("row_id", sa.BigInteger(), nullable=False),
        sa.Column("block_hash", sa.VARCHAR(length=256), nullable=False),
        sa.Column("index", sa.BigInteger(), nullable=False),
        sa.Column("path", sa.Text(), nullable=False),
        sa.Column(
            "indexed_at",
            sa.DateTime(timezone=True),
            server_default=sa.text("TIMEZONE('utc', statement_timestamp())"),
            nullable=False,
        ),
        sa.ForeignKeyConstraint(
            ["block_number"],
            ["game7_testnet_blocks.block_number"],
            name=op.f(
                "fk_game7_testnet_transactions_block_number_game7_testnet_blocks"
            ),
            ondelete="CASCADE",
        ),
        sa.PrimaryKeyConstraint("hash", name=op.f("pk_game7_testnet_transactions")),
    )
    op.create_index(
        op.f("ix_game7_testnet_transactions_block_hash"),
        "game7_testnet_transactions",
        ["block_hash"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_transactions_block_number"),
        "game7_testnet_transactions",
        ["block_number"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_transactions_from_address"),
        "game7_testnet_transactions",
        ["from_address"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_transactions_hash"),
        "game7_testnet_transactions",
        ["hash"],
        unique=True,
    )
    op.create_index(
        op.f("ix_game7_testnet_transactions_index"),
        "game7_testnet_transactions",
        ["index"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_transactions_selector"),
        "game7_testnet_transactions",
        ["selector"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_transactions_to_address"),
        "game7_testnet_transactions",
        ["to_address"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_transactions_type"),
        "game7_testnet_transactions",
        ["type"],
        unique=False,
    )
    op.create_table(
        "game7_testnet_logs",
        sa.Column("transaction_hash", sa.VARCHAR(length=256), nullable=False),
        sa.Column("block_hash", sa.VARCHAR(length=256), nullable=False),
        sa.Column("address", sa.LargeBinary(length=20), nullable=False),
        sa.Column("row_id", sa.BigInteger(), nullable=False),
        sa.Column("selector", sa.VARCHAR(length=256), nullable=True),
        sa.Column("topic1", sa.VARCHAR(length=256), nullable=True),
        sa.Column("topic2", sa.VARCHAR(length=256), nullable=True),
        sa.Column("topic3", sa.VARCHAR(length=256), nullable=True),
        sa.Column("log_index", sa.BigInteger(), nullable=False),
        sa.Column("path", sa.Text(), nullable=False),
        sa.Column(
            "indexed_at",
            sa.DateTime(timezone=True),
            server_default=sa.text("TIMEZONE('utc', statement_timestamp())"),
            nullable=False,
        ),
        sa.ForeignKeyConstraint(
            ["transaction_hash"],
            ["game7_testnet_transactions.hash"],
            name=op.f(
                "fk_game7_testnet_logs_transaction_hash_game7_testnet_transactions"
            ),
            ondelete="CASCADE",
        ),
        sa.PrimaryKeyConstraint(
            "transaction_hash", "log_index", name="pk_game7_testnet_log_index"
        ),
        sa.UniqueConstraint(
            "transaction_hash",
            "log_index",
            name="uq_game7_testnet_log_index_transaction_hash_log_index",
        ),
    )
    op.create_index(
        "idx_game7_testnet_logs_address_selector",
        "game7_testnet_logs",
        ["address", "selector"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_logs_address"),
        "game7_testnet_logs",
        ["address"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_logs_block_hash"),
        "game7_testnet_logs",
        ["block_hash"],
        unique=False,
    )
    op.create_index(
        op.f("ix_game7_testnet_logs_transaction_hash"),
        "game7_testnet_logs",
        ["transaction_hash"],
        unique=False,
    )
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_index(
        op.f("ix_game7_testnet_logs_transaction_hash"), table_name="game7_testnet_logs"
    )
    op.drop_index(
        op.f("ix_game7_testnet_logs_block_hash"), table_name="game7_testnet_logs"
    )
    op.drop_index(
        op.f("ix_game7_testnet_logs_address"), table_name="game7_testnet_logs"
    )
    op.drop_index(
        "idx_game7_testnet_logs_address_selector", table_name="game7_testnet_logs"
    )
    op.drop_table("game7_testnet_logs")
    op.drop_index(
        op.f("ix_game7_testnet_transactions_type"),
        table_name="game7_testnet_transactions",
    )
    op.drop_index(
        op.f("ix_game7_testnet_transactions_to_address"),
        table_name="game7_testnet_transactions",
    )
    op.drop_index(
        op.f("ix_game7_testnet_transactions_selector"),
        table_name="game7_testnet_transactions",
    )
    op.drop_index(
        op.f("ix_game7_testnet_transactions_index"),
        table_name="game7_testnet_transactions",
    )
    op.drop_index(
        op.f("ix_game7_testnet_transactions_hash"),
        table_name="game7_testnet_transactions",
    )
    op.drop_index(
        op.f("ix_game7_testnet_transactions_from_address"),
        table_name="game7_testnet_transactions",
    )
    op.drop_index(
        op.f("ix_game7_testnet_transactions_block_number"),
        table_name="game7_testnet_transactions",
    )
    op.drop_index(
        op.f("ix_game7_testnet_transactions_block_hash"),
        table_name="game7_testnet_transactions",
    )
    op.drop_table("game7_testnet_transactions")
    op.drop_index(
        op.f("ix_game7_testnet_reorgs_block_number"), table_name="game7_testnet_reorgs"
    )
    op.drop_index(
        op.f("ix_game7_testnet_reorgs_block_hash"), table_name="game7_testnet_reorgs"
    )
    op.drop_table("game7_testnet_reorgs")
    op.drop_index(
        op.f("ix_game7_testnet_blocks_block_timestamp"),
        table_name="game7_testnet_blocks",
    )
    op.drop_index(
        op.f("ix_game7_testnet_blocks_block_number"), table_name="game7_testnet_blocks"
    )
    op.drop_table("game7_testnet_blocks")
    # ### end Alembic commands ###