from user.schema import User

user_data = [
    User(
        user_id=1,
        username='dan',
        email='dan@dan.com',
        hashed_password='oof',
        is_admin=False,
        is_validated=True
    ),
    User(
        user_id=2,
        username='chris',
        email='chris@chris.co.uk',
        hashed_password='password',
        is_admin=False,
        is_validated=True
    ),
    User(
        user_id=3,
        username='user',
        email='user@user.com',
        hashed_password='password',
        is_admin=False,
        is_validated=True
    )
]
