from setuptools import setup

setup(name='syncmusic',
      version='1.1',
      description='Syncs music from web browser',
      author='Zack',
      author_email='zack@hypercubeplatforms.com',
      url='https://github.com/schollz/sync-music-player',
      install_requires=['eyed3', 'tornado', 'flask', 'mutagen'],
     )