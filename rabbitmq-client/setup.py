from setuptools import setup

setup(
    name='AsyncPikaClient',
    version='0.1.0',
    description='Asynchronous RabbitMQ client built on top of aio_pika',
    packages=['AsyncPikaClient'],
    install_requires=[
        'aio_pika>=9.0.0',
    ],
    classifiers=[
        'Development Status :: 3 - Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
        'Programming Language :: Python :: 3.9',
        'Programming Language :: Python :: 3.10',
    ],
)