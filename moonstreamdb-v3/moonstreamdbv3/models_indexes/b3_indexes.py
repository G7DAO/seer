# b3_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class B3BlockIndex(EvmBasedBlocks):
    __tablename__ = "b3_blocks"

class B3Reorgs(EvmBasedReorgs):
    __tablename__ = "b3_reorgs"
