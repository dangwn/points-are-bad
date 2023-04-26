PACKAGE_NAME=rabbitmq_client
VERSION=0.1.0

build:
	python setup.py sdist bdist_wheel

install:
	pip install --force-reinstall dist/$(PACKAGE_NAME)-$(VERSION)-py3-none-any.whl

clean:
	rm -rf build dist *.egg-info

