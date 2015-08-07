FROM daocloud.io/ruby:2.2.2-onbuild

MAINTAINER sllt<long@programmer.love>

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ENV LANG en_US.UTF-8
RUN mkdir -p /usr/src/app

RUN apt-get update && apt-get install --force-yes -y memcached magemagick

COPY . /usr/src/app
WORKDIR /usr/src/app
ENV USE_TAOBAO_GEM_SOURCE true

# install required gem
RUN bundle install

# compile assets
RUN bundle exec rake assets:procompile


# migrate database 
RUN bundle exec rake db:migrate

