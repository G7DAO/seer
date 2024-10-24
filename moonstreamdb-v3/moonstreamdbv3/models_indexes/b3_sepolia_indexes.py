# b3_sepolia_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class B3SepoliaBlockIndex(EvmBasedBlocks):
    __tablename__ = "b3_sepolia_blocks"

class B3SepoliaReorgs(EvmBasedReorgs):
    __tablename__ = "b3_sepolia_reorgs"
