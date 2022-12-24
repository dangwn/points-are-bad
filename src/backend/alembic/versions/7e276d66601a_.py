"""empty message

Revision ID: 7e276d66601a
Revises: 3e638a4ac37b
Create Date: 2022-12-24 12:06:56.488869

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '7e276d66601a'
down_revision = '3e638a4ac37b'
branch_labels = None
depends_on = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.create_table('matches',
    sa.Column('match_id', sa.Integer(), autoincrement=True, nullable=False),
    sa.Column('match_date', sa.String(length=10), nullable=False),
    sa.Column('home', sa.String(length=64), nullable=False),
    sa.Column('away', sa.String(length=64), nullable=False),
    sa.Column('home_goals', sa.Integer(), nullable=False),
    sa.Column('away_goals', sa.Integer(), nullable=False),
    sa.PrimaryKeyConstraint('match_id')
    )
    op.create_table('users',
    sa.Column('id', sa.Integer(), autoincrement=True, nullable=False),
    sa.Column('username', sa.String(length=30), nullable=False),
    sa.Column('email', sa.String(length=255), nullable=False),
    sa.Column('password', sa.String(length=255), nullable=False),
    sa.Column('is_admin', sa.Boolean(), nullable=False),
    sa.PrimaryKeyConstraint('id'),
    sa.UniqueConstraint('email'),
    sa.UniqueConstraint('username')
    )
    op.create_table('points',
    sa.Column('user_id', sa.Integer(), nullable=False),
    sa.Column('points', sa.Integer(), nullable=False),
    sa.Column('correct_scores', sa.Integer(), nullable=False),
    sa.Column('largest_error', sa.Integer(), nullable=False),
    sa.ForeignKeyConstraint(['user_id'], ['users.id'], ondelete='CASCADE'),
    sa.PrimaryKeyConstraint('user_id')
    )
    op.create_table('predictions',
    sa.Column('prediction_id', sa.Integer(), autoincrement=True, nullable=False),
    sa.Column('user_id', sa.Integer(), nullable=True),
    sa.Column('match_id', sa.Integer(), nullable=True),
    sa.Column('predicted_home_goals', sa.Integer(), nullable=False),
    sa.Column('predicted_away_goals', sa.Integer(), nullable=False),
    sa.Column('match_date', sa.String(length=10), nullable=False),
    sa.ForeignKeyConstraint(['match_id'], ['matches.match_id'], ondelete='CASCADE'),
    sa.ForeignKeyConstraint(['user_id'], ['users.id'], ondelete='CASCADE'),
    sa.PrimaryKeyConstraint('prediction_id')
    )
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_table('predictions')
    op.drop_table('points')
    op.drop_table('users')
    op.drop_table('matches')
    # ### end Alembic commands ###
