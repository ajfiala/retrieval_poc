"""make assistant table

Revision ID: 9e5fa6d7150a
Revises: 4cac97d41fe8
Create Date: 2024-09-21 20:38:38.763141

"""
from typing import Sequence, Union
from sqlalchemy.engine.reflection import Inspector
from alembic import op
from sqlalchemy import Column, Integer, String, DateTime, JSON, UUID
from sqlalchemy.sql import func

# revision identifiers, used by Alembic.
revision: str = '9e5fa6d7150a'
down_revision: Union[str, None] = '4cac97d41fe8'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None

def upgrade():
    # Check if app_config table already exists
    conn = op.get_bind()
    inspector = Inspector.from_engine(conn)
    tables = inspector.get_table_names()

    if 'assistant' not in tables:
        op.create_table(
            'assistant',
            Column('id', Integer,  primary_key=True, autoincrement=True),
            Column('uuid', UUID, nullable=False, unique=True),                                                                                       
            Column('name', String(255), nullable=False, unique=True),                                                                                                       
            Column('type', String(255), nullable=False),
            Column('model', String(255), nullable=False),
            Column('system_prompts', JSON, nullable=False),
            Column('metadata', JSON, nullable=True),
            Column('created_at', DateTime, server_default=func.now()),
            Column('updated_at', DateTime, server_default=func.now())

        )
        print("Table 'assistant' created successfully.")
    else:
        print("Table 'assistant' already exists.")

def downgrade():
    op.drop_table('assistant')

