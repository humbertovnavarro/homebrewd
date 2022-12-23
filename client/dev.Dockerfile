FROM node:18-bullseye
RUN apt update -y
RUN apt install wireguard-tools iproute2 iptables sudo -y

RUN echo "node  ALL = (root) NOPASSWD: /sbin/service iptables restart, /sbin/service iptables reload" >> /etc/sudoers

RUN groupadd wireguardconfig
RUN chgrp wireguardconfig /etc/wireguard
# Add "node" user from node base to group
RUN usermod -aG wireguardconfig node

WORKDIR /app
EXPOSE 3000
CMD  ["npm", "run", "dev"]

