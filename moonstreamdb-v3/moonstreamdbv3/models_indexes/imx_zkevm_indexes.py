# imx_zkevm_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class ImxZkevmBlockIndex(EvmBasedBlocks):
    __tablename__ = "imx_zkevm_blocks"

class ImxZkevmReorgs(EvmBasedReorgs):
    __tablename__ = "imx_zkevm_reorgs"
