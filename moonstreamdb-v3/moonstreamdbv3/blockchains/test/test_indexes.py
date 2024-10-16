# test_models.py

from sqlalchemy import Column, BigInteger
from models_indexes import EvmBasedBlocks, EvmBasedReorgs

class TestBlockIndex(EvmBasedBlocks):
    __tablename__ = "test_blocks"

class TestReorgs(EvmBasedReorgs):
    __tablename__ = "test_reorgs"
