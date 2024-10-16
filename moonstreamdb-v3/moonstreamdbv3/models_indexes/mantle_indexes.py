# mantle_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class MantleBlockIndex(EvmBasedBlocks):
    __tablename__ = "mantle_blocks"

class MantleReorgs(EvmBasedReorgs):
    __tablename__ = "mantle_reorgs"
