"""make message table

Revision ID: 4cac97d41fe8
Revises: 571b3155effd
Create Date: 2024-09-21 20:37:31.326852

"""
from typing import Sequence, Union
from sqlalchemy.engine.reflection import Inspector
from alembic import op
from sqlalchemy import BOOLEAN, CHAR, Column, ForeignKey, Integer, String, DateTime, JSON, UUID
from sqlalchemy.sql import func

# revision identifiers, used by Alembic.
revision: str = '4cac97d41fe8'
down_revision: Union[str, None] = '571b3155effd'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None

def upgrade():
    # Check if app_config table already exists
    conn = op.get_bind()
    inspector = Inspector.from_engine(conn)
    tables = inspector.get_table_names()

    if 'message' not in tables:
        op.create_table(
            'message',
            Column('id', Integer, primary_key=True, autoincrement=True),
            Column('uuid', UUID, nullable=False, unique=True),
            Column('user_id', UUID, ForeignKey("users.uuid"), nullable=False, unique=False),
            Column('session_id', UUID, ForeignKey("session.uuid"), nullable=False, unique=False),
            Column('user_message', String(1000), nullable=False, unique=False),
            Column('ai_message', String(1000), nullable=False, unique=False),
            Column('created_at', DateTime, server_default=func.now()),
            Column('updated_at', DateTime, server_default=func.now())
        )
        print("Table 'message' created successfully.")
    else:
        print("Table 'message' already exists.")

def downgrade():
    op.drop_table('message')