from logging.config import fileConfig
from alembic import context
from sqlalchemy import engine_from_config, pool
import os
import sys
import pkgutil
import importlib
from sqlalchemy.ext.declarative import DeclarativeMeta

# Get the Alembic Config object
config = context.config

# Interpret the config file for Python logging
if config.config_file_name is not None:
    fileConfig(config.config_file_name)

# Add project directory to sys.path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

# Import Base and target metadata
from models_indexes.abstract_indexes import Base as IndexesBase

target_metadata = IndexesBase.metadata


# Dynamically import all models
def import_models(package_name):
    package = importlib.import_module(package_name)
    ###print(dir(package))
    for loader, module_name, is_pkg in pkgutil.walk_packages(
        package.__path__, package_name + "."
    ):
        importlib.import_module(module_name)


import_models("moonstreamdbv3.models_indexes")


# Collect all model classes
def get_model_classes(base):
    return base.__subclasses__() + [
        cls for scls in base.__subclasses__() for cls in get_model_classes(scls)
    ]


model_classes = get_model_classes(IndexesBase)

# Generate set of table names
model_tablenames = {
    cls.__tablename__ for cls in model_classes if hasattr(cls, "__tablename__")
}

print(model_tablenames)


# Define include_symbol function
def include_symbol(tablename, schema):
    return tablename in model_tablenames


# Run migrations offline
def run_migrations_offline():
    url = config.get_main_option("sqlalchemy.url")
    context.configure(
        url=url,
        target_metadata=target_metadata,
        include_symbol=include_symbol,
        literal_binds=True,
        dialect_opts={"paramstyle": "named"},
        version_table="alembic_version",
        include_schemas=True,
    )

    with context.begin_transaction():
        context.run_migrations()


# Run migrations online
def run_migrations_online():
    connectable = engine_from_config(
        config.get_section(config.config_ini_section, {}),
        prefix="sqlalchemy.",
        poolclass=pool.NullPool,
    )

    with connectable.connect() as connection:
        context.configure(
            connection=connection,
            target_metadata=target_metadata,
            include_symbol=include_symbol,
            version_table="alembic_version",
            include_schemas=True,
        )

        with context.begin_transaction():
            context.run_migrations()


# Determine if running offline or online
if context.is_offline_mode():
    run_migrations_offline()
else:
    run_migrations_online()
