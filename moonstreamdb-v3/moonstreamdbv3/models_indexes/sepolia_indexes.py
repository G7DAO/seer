# sepolia_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class SepoliaBlockIndex(EvmBasedBlocks):
    __tablename__ = "sepolia_blocks"

class SepoliaReorgs(EvmBasedReorgs):
    __tablename__ = "sepolia_reorgs"
