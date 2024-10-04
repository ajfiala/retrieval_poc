"""make users

Revision ID: 6ab78259c9db
Revises: 
Create Date: 2024-09-21 20:35:41.729894

"""
from typing import Sequence, Union
from sqlalchemy.engine.reflection import Inspector

from alembic import op
from sqlalchemy import BOOLEAN, Column, Integer, String, DateTime, UUID
from sqlalchemy.sql import func



# revision identifiers, used by Alembic.
revision: str = '6ab78259c9db'
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade():
    # Check if app_config table already exists
    conn = op.get_bind()
    inspector = Inspector.from_engine(conn)
    tables = inspector.get_table_names()

    if 'users' not in tables:
        op.create_table(
            'users',
            Column('id', Integer, primary_key=True, autoincrement=True),
            Column('uuid', UUID, nullable=False, unique=True),
            Column('name', String(100), nullable=False, unique=False),
            Column('active', BOOLEAN, nullable=False, default=False),
            Column('created_at', DateTime, server_default=func.now()),
            Column('updated_at', DateTime, server_default=func.now())
        )
        print("Table 'users' created successfully.")
    else:
        print("Table 'users' already exists.")

def downgrade():
    op.drop_table('users')
