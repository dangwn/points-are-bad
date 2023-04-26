from setuptools import setup

setup(
    name='rabbitmq_client',
    version='0.1.0',
    description='Asynchronous RabbitMQ client built on top of aio_pika',
    packages=['rabbitmq_client'],
    install_requires=[
        'aio_pika>=9.0.0',
    ],
    classifiers=[
        'Development Status :: 3 - Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.8',
        'Programming Language :: Python :: 3.9',
        'Programming Language :: Python :: 3.10',
    ],
)