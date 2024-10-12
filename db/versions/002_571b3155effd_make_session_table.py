"""make session table

Revision ID: 571b3155effd
Revises: 6ab78259c9db
Create Date: 2024-09-21 20:36:54.598776

"""
from typing import Sequence, Union
from sqlalchemy.engine.reflection import Inspector
from alembic import op
from sqlalchemy import BOOLEAN, CHAR, Column, ForeignKey, Integer, String, DateTime, UUID
from sqlalchemy.sql import func

# revision identifiers, used by Alembic.
revision: str = '571b3155effd'
down_revision: Union[str, None] = '6ab78259c9db'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None

def upgrade():
    # Check if app_config table already exists
    conn = op.get_bind()
    inspector = Inspector.from_engine(conn)
    tables = inspector.get_table_names()

    if 'session' not in tables:
        op.create_table(
            'session',
            Column('id', Integer, primary_key=True, autoincrement=True),
            Column('uuid', UUID, nullable=False, unique=True),
            Column('user_id', UUID, ForeignKey("users.uuid"), nullable=False, unique=False),
            Column('active', BOOLEAN, nullable=False, default=True),
            Column('created_at', DateTime, server_default=func.now()),
            Column('updated_at', DateTime, server_default=func.now())
        )
        print("Table 'session' created successfully.")
    else:
        print("Table 'session' already exists.")

def downgrade():
    op.drop_table('session')
