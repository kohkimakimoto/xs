FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
STOPSIGNAL SIGRTMIN+3

RUN rm /etc/dpkg/dpkg.cfg.d/excludes && apt-get update && apt-get install -y \
  systemd sudo openssh-server vim locales

# Setup locale
RUN echo 'LANG="en_US.UTF-8"' > /etc/default/locale

# Setup sshd
RUN sed -i 's/#Port 22/Port 22/' /etc/ssh/sshd_config && \
  sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config && \
  sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config

# Setup dev user account
ARG DEV_USER=xs-test-user
RUN useradd -m -s /bin/bash ${DEV_USER} && \
  mkdir /home/${DEV_USER}/.ssh && \
  chmod 700 /home/${DEV_USER}/.ssh && \
  echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDAZ5UYwwIA5a9HNhdOjHX8O+5ZtydhF1qn07d4k+8E8X/K4VS4DwBueAK2F4e9qONcme8VZYv/BPep1EK0QP2hezKucTrRExwcc6idQog5bbzFq3j7AsLPWQap2Dil5dbu2NyHO5nIjS0mNRRfYK2Q+1S5SVJnqUZZVmNle8T0TM6JyH/ahM3REoVybX1dsEbqYPavsuALOeoYZJq/AME91K/a5lL2zlqApbDQUzMpV5YZMk9DFziMETwkkRN/b3c5MYHde205e3lfVgpFtYa4E+XnWx2WHz0+DyWjXfFlLHi5ajcNzpu/iS7JoPb61JGsRyUbNKTW4P3qX3Gbr5rNIUG7+hlJBirup5lOjjh5ukgi51cVxzGXzlGFyL0E9mLxXUIbZ7BPMHt/PUMmmdRHwsoiay3s/qzWiXY6jwSRPrhlROWWK5QddtwrKm5g5DsoPfZhQP6+5KwBusgtkbrzAk95BhhYukgtAvh/RIIbQ2+IKRFZLvlw3dAQ0KHRr3EGW4lcepJBqzKlGs2mgwGsYkWp4n4BarsWJmg3bnJQNkACdFuGtsCUiUAmGlUBg4/on9M8G5NVQdF73pRN4ownxN9sybHlKXsveYpqjKh+scM9BV6dsDJMu0Nf39qOD3ibERuWRm764qCAveYb6W1iNuLWcseKtN9QMPWjalg0qQ== ${DEV_USER}" > /home/${DEV_USER}/.ssh/authorized_keys && \
  chmod 600 /home/${DEV_USER}/.ssh/authorized_keys && \
  chown -R ${DEV_USER}:${DEV_USER} /home/${DEV_USER}/.ssh && \
  echo "${DEV_USER} ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

EXPOSE 22
CMD ["/sbin/init"]
