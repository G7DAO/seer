"""
Moonstream database V3

Example of label_data column record:
    {
        "label": "ERC20",
        "label_data": {
            "name": "Uniswap",
            "symbol": "UNI"
        }
    },
    {
        "label": "Exchange"
        "label_data": {...}
    }
"""

import uuid

from sqlalchemy import (
    VARCHAR,
    BigInteger,
    Column,
    DateTime,
    Index,
    Integer,
    LargeBinary,
    MetaData,
    Text,
)
from sqlalchemy.dialects.postgresql import JSONB, UUID
from sqlalchemy.ext.compiler import compiles
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.sql import expression, text

"""
Naming conventions doc
https://docs.sqlalchemy.org/en/20/core/constraints.html#configuring-constraint-naming-conventions
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

"""
Creating a utcnow function which runs on the Posgres database server when created_at and updated_at
fields are populated.
Following:
1. https://docs.sqlalchemy.org/en/13/core/compiler.html#utc-timestamp-function
2. https://www.postgresql.org/docs/current/functions-datetime.html#FUNCTIONS-DATETIME-CURRENT
3. https://stackoverflow.com/a/33532154/13659585
"""


class utcnow(expression.FunctionElement):
    type = DateTime  # type: ignore


@compiles(utcnow, "postgresql")
def pg_utcnow(element, compiler, **kwargs):
    return "TIMEZONE('utc', statement_timestamp())"


class EvmBasedLabel(Base):  # type: ignore
    __abstract__ = True

    id = Column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid.uuid4,
        unique=True,
        nullable=False,
    )
    label = Column(VARCHAR(256), nullable=False, index=True)

    transaction_hash = Column(
        VARCHAR(128),
        nullable=False,
        index=True,
    )
    log_index = Column(Integer, nullable=True)

    block_number = Column(
        BigInteger,
        nullable=False,
        index=True,
    )
    block_hash = Column(VARCHAR(256), nullable=False)
    block_timestamp = Column(BigInteger, nullable=False)

    caller_address = Column(
        LargeBinary,
        nullable=True,
        index=True,
    )
    origin_address = Column(
        LargeBinary,
        nullable=True,
        index=True,
    )

    address = Column(
        LargeBinary,
        nullable=False,
        index=True,
    )

    label_name = Column(Text, nullable=True, index=True)
    label_type = Column(VARCHAR(64), nullable=True, index=True)
    label_data = Column(JSONB, nullable=True)

    created_at = Column(
        DateTime(timezone=True), server_default=utcnow(), nullable=False
    )
