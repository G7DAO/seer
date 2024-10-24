# xai_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class XaiBlockIndex(EvmBasedBlocks):
    __tablename__ = "xai_blocks"
    l1_block_number = Column(BigInteger, nullable=False)

class XaiReorgs(EvmBasedReorgs):
    __tablename__ = "xai_reorgs"
