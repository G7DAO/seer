# mantle_sepolia_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class MantleSepoliaBlockIndex(EvmBasedBlocks):
    __tablename__ = "mantle_sepolia_blocks"

class MantleSepoliaReorgs(EvmBasedReorgs):
    __tablename__ = "mantle_sepolia_reorgs"
