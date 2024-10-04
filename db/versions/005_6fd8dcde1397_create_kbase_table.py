"""create kbase table

Revision ID: 6fd8dcde1397
Revises: 9e5fa6d7150a
Create Date: 2024-09-22 15:05:45.791203

"""
from typing import Sequence, Union
from sqlalchemy.engine.reflection import Inspector
from alembic import op
from sqlalchemy import Column, Integer, String, DateTime, JSON, UUID, ForeignKey
from sqlalchemy.sql import func


# revision identifiers, used by Alembic.
revision: str = '6fd8dcde1397'
down_revision: Union[str, None] = '9e5fa6d7150a'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None

def upgrade():
    conn = op.get_bind()
    inspector = Inspector.from_engine(conn)

    # Check if the 'app_config' table exists
    if 'kbase' not in inspector.get_table_names():
        op.create_table(
            'kbase',
            Column('id', Integer, primary_key=True, autoincrement=True),
            Column('uuid', UUID, unique=True, nullable=False),
            Column('name', String(255), nullable=False, unique=True),
            Column('description', String(2500), nullable=True),
            Column('created_at', DateTime, server_default=func.now()),
            Column('updated_at', DateTime, server_default=func.now())
        )
    else:
        print("Table 'kbase' already exists.")

def downgrade():
    op.drop_table('kbase')