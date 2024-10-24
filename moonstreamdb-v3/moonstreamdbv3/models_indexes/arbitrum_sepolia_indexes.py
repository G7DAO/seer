# arbitrum_sepolia_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class ArbitrumSepoliaBlockIndex(EvmBasedBlocks):
    __tablename__ = "arbitrum_sepolia_blocks"
    l1_block_number = Column(BigInteger, nullable=False)

class ArbitrumSepoliaReorgs(EvmBasedReorgs):
    __tablename__ = "arbitrum_sepolia_reorgs"
