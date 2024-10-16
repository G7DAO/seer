# arbitrum_one_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class ArbitrumOneBlockIndex(EvmBasedBlocks):
    __tablename__ = "arbitrum_one_blocks"
    l1_block_number = Column(BigInteger, nullable=False)

class ArbitrumOneReorgs(EvmBasedReorgs):
    __tablename__ = "arbitrum_one_reorgs"
