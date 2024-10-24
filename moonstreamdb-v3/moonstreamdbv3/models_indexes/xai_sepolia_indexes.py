# xai_sepolia_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class XaiSepoliaBlockIndex(EvmBasedBlocks):
    __tablename__ = "xai_sepolia_blocks"
    l1_block_number = Column(BigInteger, nullable=False)

class XaiSepoliaReorgs(EvmBasedReorgs):
    __tablename__ = "xai_sepolia_reorgs"
