"""empty message

Revision ID: 3e638a4ac37b
Revises: bbaf99c20112
Create Date: 2022-12-22 09:23:41.950445

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '3e638a4ac37b'
down_revision = 'bbaf99c20112'
branch_labels = None
depends_on = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.alter_column('predictions', 'match_date',
               existing_type=sa.VARCHAR(length=10),
               nullable=False)
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.alter_column('predictions', 'match_date',
               existing_type=sa.VARCHAR(length=10),
               nullable=True)
    # ### end Alembic commands ###