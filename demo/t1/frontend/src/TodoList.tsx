import React, { useState } from 'react';
import { Button, Checkbox, Input, List, Space } from 'antd';
import {
    DeleteOutlined,
    DownloadOutlined,
    EditOutlined,
    ExportOutlined,
    EyeOutlined,
    ImportOutlined,
    LoginOutlined,
    LogoutOutlined,
    RobotOutlined,
    ShareAltOutlined,
    UploadOutlined
} from '@ant-design/icons';

import {Aid} from "aid-js-sdk";

interface Todo {
    id: number;
    task: string;
    done: boolean;
}

interface ActionButton {
    icon: React.ReactNode;
    text: string;
    callback?: () => void;
}

export const TodoList: React.FC = () => {
    const [todos, setTodos] = useState<Todo[]>([]);
    const [inputValue, setInputValue] = useState<string>('');
    const [editingId, setEditingId] = useState<number | null>(null);

    const addTodo = (): void => {
        if (inputValue.trim() !== '') {
            setTodos([...todos, { id: Date.now(), task: inputValue, done: false }]);
            setInputValue('');
        }
    };

    const deleteTodo = (id: number): void => {
        setTodos(todos.filter(todo => todo.id !== id));
    };

    const toggleComplete = (id: number): void => {
        setTodos(todos.map(todo =>
            todo.id === id ? { ...todo, done: !todo.done } : todo
        ));
    };

    const startEditing = (id: number): void => {
        setEditingId(id);
    };

    const finishEditing = (id: number, newText: string): void => {
        setTodos(todos.map(todo =>
            todo.id === id ? { ...todo, task: newText } : todo
        ));
    };

    const actionButtons: ActionButton[] = [
        { icon: <LoginOutlined />, text: 'Login' },
        { icon: <LogoutOutlined />, text: 'Logout' },
        { icon: <UploadOutlined />, text: 'Upload' },
        { icon: <DownloadOutlined />, text: 'Download' },
        { icon: <ShareAltOutlined />, text: 'Share' },
        { icon: <EyeOutlined />, text: 'View' },
        { icon: <RobotOutlined />, text: 'Generate Aid', callback: () => {
            const uuid = crypto.randomUUID();
            const aid = new Aid(uuid, new Map(), new Map(), new Map());
            alert(aid.toStr());
        } },
        { icon: <ExportOutlined />, text: 'Export Aid' },
        { icon: <ImportOutlined />, text: 'Import Aid' },
    ];

    return (
        <div style={{ maxWidth: 600, margin: '0 auto', padding: 20 }}>
            <h1>TodoList</h1>
            <Space wrap style={{ marginBottom: 16 }}>
                {actionButtons.map((button, index) => (
                    <Button key={index} icon={button.icon} onClick={button.callback}>
                        {button.text}
                    </Button>
                ))}
            </Space>
            <Space.Compact style={{ width: '100%' }}>
                <Input
                    value={inputValue}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setInputValue(e.target.value)}
                    onPressEnter={addTodo}
                    placeholder="Add a new todo"
                />
                <Button type="primary" onClick={addTodo}>Add</Button>
            </Space.Compact>
            <List
                style={{ marginTop: 20 }}
                bordered
                dataSource={todos}
                renderItem={(todo: Todo) => (
                    <List.Item>
                        <Checkbox
                            checked={todo.done}
                            onChange={() => toggleComplete(todo.id)}
                        />
                        {editingId === todo.id ? (
                            <Input
                                value={todo.task}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) => finishEditing(todo.id, e.target.value)}
                                onBlur={() => setEditingId(null)}
                                autoFocus
                            />
                        ) : (
                            <span style={{
                                textDecoration: todo.done ? 'line-through' : 'none',
                                marginLeft: 8,
                                marginRight: 8,
                                flex: 1
                            }}>
                                {todo.task}
                            </span>
                        )}
                        <Space>
                            <Button
                                icon={<EditOutlined />}
                                onClick={() => startEditing(todo.id)}
                            />
                            <Button
                                icon={<DeleteOutlined />}
                                onClick={() => deleteTodo(todo.id)}
                                danger
                            />
                        </Space>
                    </List.Item>
                )}
            />
        </div>
    );
};