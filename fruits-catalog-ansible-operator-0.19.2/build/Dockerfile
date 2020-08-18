FROM quay.io/operator-framework/ansible-operator:v0.19.2

COPY requirements.yml ${HOME}/requirements.yml
RUN ansible-galaxy collection install -r ${HOME}/requirements.yml \
 && chmod -R ug+rwx ${HOME}/.ansible

COPY watches.yaml ${HOME}/watches.yaml
COPY ansible.cfg /etc/ansible/ansible.cfg

COPY roles/ ${HOME}/roles/
