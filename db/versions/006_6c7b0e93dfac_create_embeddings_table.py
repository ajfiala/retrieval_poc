"""create embeddings table

Revision ID: 6c7b0e93dfac
Revises: 6fd8dcde1397
Create Date: 2024-09-22 15:18:03.279277

"""
from typing import Sequence, Union
from sqlalchemy.engine.reflection import Inspector
from alembic import op
from sqlalchemy import Column, Integer, String, DateTime, JSON, UUID, ForeignKey, FLOAT, ARRAY
from sqlalchemy.sql import func
from pgvector.sqlalchemy import Vector

# revision identifiers, used by Alembic.
revision: str = '6c7b0e93dfac'
down_revision: Union[str, None] = '6fd8dcde1397'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None

def upgrade():
    conn = op.get_bind()
    inspector = Inspector.from_engine(conn)

    # enable PGVector
    op.execute("CREATE EXTENSION IF NOT EXISTS vector;")

    # Check if the 'app_config' table exists
    if 'kbase_embeddings' not in inspector.get_table_names():
        op.create_table(
            'kbase_embeddings',
            Column('id', Integer, primary_key=True, autoincrement=True),
            Column('uuid', UUID, unique=True, nullable=False),
            Column('kbase_id', UUID, ForeignKey('kbase.uuid'), nullable=False),
            Column('chunk_id', Integer, nullable=False),
            Column('content', String(5000), nullable=False),
            Column('embedding', Vector, nullable=False),
            Column('metadata', JSON, nullable=True),
            Column('created_at', DateTime, server_default=func.now()),
            Column('updated_at', DateTime, server_default=func.now())
        )
    else:
        print("Table 'kbase' already exists.")

def downgrade():
    op.drop_table('kbase_embeddings')