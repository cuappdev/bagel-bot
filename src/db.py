from enum import Enum, unique

import sqlalchemy as db
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, relationship

import slack


Engine = db.create_engine('sqlite:///bagel.sqlite') 
session = sessionmaker(bind=Engine)()


Base = declarative_base()


user_chat_association_table = db.Table(
    'user_chat_association_table',
    Base.metadata,
    db.Column('user_id', db.Integer, db.ForeignKey('users.id')),
    db.Column('chat_id', db.Integer, db.ForeignKey('chats.id'))
)


"""
A bagel instance represents the instance when groups are formed in the bagel chat.

It is used to send reminders for group chats. 
"""
class BagelInstance(Base):
    __tablename__ = 'bagel_instance'
    id = db.Column(db.Integer, primary_key=True) 
    bagel_date = db.Column(db.Integer)
    chats = relationship('Chat')

    @staticmethod
    def query():
        return session.query(BagelInstance)


    def __init__(self, bagel_date):
        self.bagel_date = bagel_date
        self.chats = []


class Chat(Base):
    __tablename__ = 'chat'
    id = db.Column(db.Integer, primary_key=True)
    slack_id = db.Column(db.Integer)
    users = relationship(
        'User',
        secondary=user_chat_association_table,
        back_populates='chats')
    bagel_instance_id = db.Column(db.Integer, db.ForeignKey('bagel_instance.id'))
    status = db.Column(db.Integer)

    @staticmethod
    def query():
        return session.query(Chat)

    def __init__(self, slack_id):
        self.slack_id = slack_id
        self.users = []
        self.status = ChatStatus.UNSCHEDULED.value

    def set_status(self, status):
        self.status = status.value
        session.commit()

    def get_status(self):
        return ChatStatus.from_value(self.status)


@unique
class ChatStatus(Enum):
    UNSCHEDULED = 0 
    SCHEDULED = 1
    MET = 2

    @staticmethod
    def from_value(raw_value):
        if raw_value == 0:
            return ChatStatus.UNSCHEDULED
        elif raw_value == 1:
            return ChatStatus.SCHEDULED
        elif raw_value == 2:
            return ChatStatus.MET
        else:
            return ChatStatus.UNSCHEDULED


class User(Base): 
    __tablename__ = 'user'
    id = db.Column(db.Integer, primary_key=True)
    slack_id = db.Column(db.Integer)
    name = db.Column(db.String)
    chats = relationship(
        'Chat', 
        secondary=user_chat_association_table,
        back_populates='users')
    is_active = db.Column(db.Boolean)

    @staticmethod
    def query():
        return session.query(User)
    
    def __init__(self, slack_id, name, is_active):
        self.slack_id = slack_id
        self.name = name
        self.is_active = is_active


Base.metadata.create_all(Engine)
