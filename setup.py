import setuptools

requires = [
    'Sphinx >= 1.5',
    'six',
]


def readme():
    try:
        with open('README.rst') as f:
            return f.read()
    except IOError:
        pass


setuptools.setup(
    name='Elephant.Node',
    version='0.0.1',
    url='https://github.com/elastos/Elastos.ELA.Elephant.Node',
    license='MIT',
    author='huangxiaoming',
    author_email='huangxiaoming@elastos.org',
    description='elastos blockchain entry port',
    long_description=readme(),
    zip_safe=False,
    classifiers=[
        'Development Status :: 5 - Production/Stable',
        'Environment :: Console',
        'Environment :: Web Environment',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Operating System :: OS Independent',
        'Programming Language :: Python',
        'Programming Language :: Python :: 2.7',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
        'Topic :: Documentation',
        'Topic :: Utilities',
    ],
    platforms='any',
    packages=setuptools.find_packages(),
    include_package_data=True,
    install_requires=requires,
)