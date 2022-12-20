"""empty message

Revision ID: 7d3ac0c989f4
Revises: 5c71dbd39513
Create Date: 2022-12-17 10:30:29.362727

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '7d3ac0c989f4'
down_revision = '5c71dbd39513'
branch_labels = None
depends_on = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_column('points', 'prev_day_points')
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.add_column('points', sa.Column('prev_day_points', sa.INTEGER(), autoincrement=False, nullable=False))
    # ### end Alembic commands ###
