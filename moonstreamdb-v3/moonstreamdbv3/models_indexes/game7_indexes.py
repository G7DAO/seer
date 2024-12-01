# game7_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class Game7BlockIndex(EvmBasedBlocks):
    __tablename__ = "game7_blocks"

    l1_block_number = Column(BigInteger, nullable=False)


class Game7Reorgs(EvmBasedReorgs):
    __tablename__ = "game7_reorgs"
