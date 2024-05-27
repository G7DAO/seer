import os
from contextlib import contextmanager
from typing import Generator

from sqlalchemy import create_engine
from sqlalchemy.orm import Session, sessionmaker


MOONSTREAM_DB_URI = os.environ.get("MOONSTREAM_INDEX_URI")
if MOONSTREAM_DB_URI is None:
    raise ValueError("MOONSTREAM_DB_URI environment variable must be set")

engine = create_engine(MOONSTREAM_DB_URI)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def yield_db_session() -> Generator[Session, None, None]:
    """
    Yields a database connection (created using environment variables).
    As per FastAPI docs:
    https://fastapi.tiangolo.com/tutorial/sql-databases/#create-a-dependency
    """
    session = SessionLocal()
    try:
        yield session
    finally:
        session.close()


yield_db_session_ctx = contextmanager(yield_db_session)
