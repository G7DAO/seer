# imx_zkevm_sepolia_models.py

from sqlalchemy import Column, BigInteger
from models_indexes.abstract_indexes import EvmBasedBlocks, EvmBasedReorgs

class ImxZkevmSepoliaBlockIndex(EvmBasedBlocks):
    __tablename__ = "imx_zkevm_sepolia_blocks"

class ImxZkevmSepoliaReorgs(EvmBasedReorgs):
    __tablename__ = "imx_zkevm_sepolia_reorgs"
