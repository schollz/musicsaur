from setuptools import setup

setup(name='syncmusic',
      version='1.2',
      description='Syncs music from web browser',
      author='Zack',
      author_email='zack@hypercubeplatforms.com',
      url='https://github.com/schollz/musicsaur',
      install_requires=['tornado', 'mutagen','configparser'],
     )