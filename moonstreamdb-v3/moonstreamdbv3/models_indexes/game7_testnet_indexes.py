# game7_testnet_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class Game7TestnetBlockIndex(EvmBasedBlocks):
    __tablename__ = "game7_testnet_blocks"
    l1_block_number = Column(BigInteger, nullable=False)

class Game7TestnetReorgs(EvmBasedReorgs):
    __tablename__ = "game7_testnet_reorgs"
