import uuid

from sqlalchemy import (
    VARCHAR,
    BigInteger,
    Boolean,
    Column,
    DateTime,
    ForeignKey,
    Index,
    Integer,
    LargeBinary,
    MetaData,
    PrimaryKeyConstraint,
    Text,
    UniqueConstraint,
)
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.ext.compiler import compiles
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.sql import expression

"""
Naming conventions doc
https://docs.sqlalchemy.org/en/13/core/constraints.html#configuring-constraint-naming-conventions
"""
convention = {
    "ix": "ix_%(column_0_label)s",
    "uq": "uq_%(table_name)s_%(column_0_name)s",
    "ck": "ck_%(table_name)s_%(constraint_name)s",
    "fk": "fk_%(table_name)s_%(column_0_name)s_%(referred_table_name)s",
    "pk": "pk_%(table_name)s",
}
metadata = MetaData(naming_convention=convention)
Base = declarative_base(metadata=metadata)


class utcnow(expression.FunctionElement):
    type = DateTime  # type: ignore


@compiles(utcnow, "postgresql")
def pg_utcnow(element, compiler, **kwargs):
    return "TIMEZONE('utc', statement_timestamp())"


class EvmBasedBlocks(Base):
    __abstract__ = True

    block_number = Column(BigInteger, primary_key=True, nullable=False, index=True)
    block_hash = Column(VARCHAR(256), nullable=False, index=False)
    block_timestamp = Column(BigInteger, nullable=False, index=True)
    parent_hash = Column(VARCHAR(256), nullable=False)
    row_id = Column(BigInteger, nullable=False, index=False)
    path = Column(Text, nullable=False)
    transactions_indexed_at = Column(DateTime(timezone=True), nullable=True)
    logs_indexed_at = Column(DateTime(timezone=True), nullable=True)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )


class EvmBasedTransactions(Base):
    __abstract__ = True

    hash = Column(
        VARCHAR(256), primary_key=True, unique=True, nullable=False, index=True
    )
    from_address = Column(LargeBinary(length=20), nullable=False, index=True)
    to_address = Column(LargeBinary(length=20), nullable=False, index=True)
    selector = Column(VARCHAR(256), nullable=True, index=True)
    type = Column(Integer, nullable=True, index=True)
    row_id = Column(BigInteger, nullable=False, index=False)
    block_hash = Column(VARCHAR(256), nullable=False, index=True)
    index = Column(BigInteger, nullable=False, index=True)
    path = Column(Text, nullable=False)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )


class EvmBasedLogs(Base):
    __abstract__ = True

    block_hash = Column(VARCHAR(256), nullable=False, index=True)
    address = Column(LargeBinary(length=20), nullable=False, index=True)
    row_id = Column(BigInteger, nullable=False, index=False)
    selector = Column(VARCHAR(256), nullable=True, index=False)
    topic1 = Column(VARCHAR(256), nullable=True, index=False)
    topic2 = Column(VARCHAR(256), nullable=True, index=False)
    topic3 = Column(VARCHAR(256), nullable=True, index=False)
    log_index = Column(BigInteger, nullable=False, index=False)
    path = Column(Text, nullable=False)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )


class EvmBasedReorgs(Base):
    __abstract__ = True
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    block_number = Column(BigInteger, nullable=False, index=True)
    block_hash = Column(VARCHAR(256), nullable=False, index=True)


class AbiJobs(Base):
    __tablename__ = "abi_jobs"

    __table_args__ = (
        UniqueConstraint(
            "chain",
            "address",
            "abi_selector",
            "customer_id",
            name="uq_abi_jobs",
        ),
    )

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    address = Column(LargeBinary, nullable=False, index=True)
    user_id = Column(UUID(as_uuid=True), nullable=False, index=True)
    customer_id = Column(UUID(as_uuid=True), nullable=True, index=True)
    abi_selector = Column(VARCHAR(256), nullable=False, index=True)
    chain = Column(VARCHAR(256), nullable=False, index=True)
    abi_name = Column(VARCHAR(256), nullable=False, index=True)
    status = Column(VARCHAR(256), nullable=False, index=True)
    historical_crawl_status = Column(VARCHAR(256), nullable=False, index=True)
    progress = Column(Integer, nullable=False, index=False)
    moonworm_task_pickedup = Column(Boolean, nullable=False, index=False)
    deployment_block_number = Column(BigInteger, nullable=True, index=False)
    abi = Column(Text, nullable=False)
    created_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )
    updated_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )


class AbiSubscriptions(Base):
    __tablename__ = "abi_subscriptions"

    __table_args__ = (PrimaryKeyConstraint("abi_job_id", "subscription_id"),)

    abi_job_id = Column(
        UUID(as_uuid=True),
        ForeignKey("abi_jobs.id"),
        nullable=False,
        index=True,
    )
    subscription_id = Column(UUID(as_uuid=True), nullable=False, index=True)
    created_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )