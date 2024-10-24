# polygon_models.py

from sqlalchemy import Column, BigInteger
from models_indexes import EvmBasedBlocks, EvmBasedReorgs

class PolygonBlockIndex(EvmBasedBlocks):
    __tablename__ = "polygon_blocks"

class PolygonReorgs(EvmBasedReorgs):
    __tablename__ = "polygon_reorgs"
