# ethereum_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class EthereumBlockIndex(EvmBasedBlocks):
    __tablename__ = "ethereum_blocks"

class EthereumReorgs(EvmBasedReorgs):
    __tablename__ = "ethereum_reorgs"
