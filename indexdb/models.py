import uuid

from sqlalchemy import (
    VARCHAR,
    BigInteger,
    Column,
    DateTime,
    Integer,
    MetaData,
    Text,
    Boolean
)
from sqlalchemy.dialects.postgresql import JSONB, UUID
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


class EthereumBlockIndex(Base):
    __tablename__ = "ethereum_index"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    block_number = Column(BigInteger, nullable=False, index=True)
    block_hash = Column(VARCHAR(256), nullable=False, index=True)
    block_timestamp = Column(DateTime, nullable=False, index=True)
    parent_hash = Column(VARCHAR(256), nullable=False)
    path = Column(Text, nullable=False)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )

class EthereumTransactionIndex(Base):
    __tablename__ = "ethereum_transaction_index"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    block_number = Column(BigInteger, nullable=False, index=True)
    block_hash = Column(VARCHAR(256), nullable=False, index=True)
    block_timestamp = Column(DateTime, nullable=False, index=True)
    transaction_hash = Column(VARCHAR(256), nullable=False, index=True)
    transaction_index = Column(BigInteger, nullable=False, index=True)
    path = Column(Text, nullable=False)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )

class EthereumLogIndex(Base):
    __tablename__ = "ethereum_log_index"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    block_number = Column(BigInteger, nullable=False, index=True)
    block_hash = Column(VARCHAR(256), nullable=False, index=True)
    block_timestamp = Column(DateTime, nullable=False, index=True)
    transaction_hash = Column(VARCHAR(256), nullable=False, index=True)
    selector = Column(VARCHAR(256), nullable=False, index=True)
    topic1 = Column(VARCHAR(256), nullable=False, index=True)
    topic2 = Column(VARCHAR(256), nullable=False, index=True)
    log_index = Column(BigInteger, nullable=False, index=True)
    path = Column(Text, nullable=False)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )


class PolygonBlockIndex(Base):
    __tablename__ = "polygon_index"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    block_number = Column(BigInteger, nullable=False, index=True)
    block_hash = Column(VARCHAR(256), nullable=False, index=True)
    block_timestamp = Column(DateTime, nullable=False, index=True)
    parent_hash = Column(VARCHAR(256), nullable=False, index=True)
    path = Column(Text, nullable=False)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )

class PolygonTransactionIndex(Base):
    __tablename__ = "polygon_transaction_index"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    block_number = Column(BigInteger, nullable=False, index=True)
    block_hash = Column(VARCHAR(256), nullable=False, index=True)
    block_timestamp = Column(DateTime, nullable=False, index=True)
    transaction_hash = Column(VARCHAR(256), nullable=False, index=True)
    transaction_index = Column(BigInteger, nullable=False, index=True)
    path = Column(Text, nullable=False)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )

class PolygonLogIndex(Base):
    __tablename__ = "polygon_log_index"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    block_number = Column(BigInteger, nullable=False, index=True)
    block_hash = Column(VARCHAR(256), nullable=False, index=True)
    block_timestamp = Column(DateTime, nullable=False, index=True)
    transaction_hash = Column(VARCHAR(256), nullable=False, index=True)
    selector = Column(VARCHAR(256), nullable=False, index=True)
    topic1 = Column(VARCHAR(256), nullable=False, index=True)
    topic2 = Column(VARCHAR(256), nullable=False, index=True)
    log_index = Column(BigInteger, nullable=False, index=False)
    path = Column(Text, nullable=False)
    indexed_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )

class AbiJobs(Base):
    __tablename__ = "abi_jobs"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    address = Column(VARCHAR(256), nullable=False, index=True)
    user_id = Column(UUID(as_uuid=True), nullable=False, index=True)
    customer_id = Column(UUID(as_uuid=True), nullable=True, index=True)
    abi_selector = Column(VARCHAR(256), nullable=False, index=True)
    chain = Column(VARCHAR(256), nullable=False, index=True)
    abi_name = Column(VARCHAR(256), nullable=False, index=True)
    status = Column(VARCHAR(256), nullable=False, index=True)
    historical_crawl_status = Column(VARCHAR(256), nullable=False, index=True)
    progress = Column(Integer, nullable=False, index=False)
    moonworm_task_pickedup = Column(Boolean, nullable=False, index=False)
    abi = Column(Text, nullable=False)
    created_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )
    updated_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )


