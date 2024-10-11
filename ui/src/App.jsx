import React, { useState, useRef, useEffect } from 'react';
import useChat from './lib/store';
import { useShallow } from 'zustand/react/shallow';

// ChatMessage component for individual messages
const ChatMessage = ({ message, isUser }) => (
    <div style={{ marginBottom: '10px', textAlign: isUser ? 'right' : 'left' }}>
        <div
            style={{
                display: 'inline-block',
                padding: '10px',
                backgroundColor: isUser ? '#c0c0c0' : '#ffffff',
                color: '#000000',
                border: '1px solid #000',
                fontFamily: 'Arial, sans-serif',
            }}
        >
            {message}
        </div>
    </div>
);

// Main App component
export default function App() {
    const [input, setInput, messages, sendMessage, isThinking] = useChat(useShallow(state => [state.input, state.setInput, state.messages, state.sendMessage, state.isThinking]))

    const chatContainerRef = useRef(null);

    // Function to send message to /chat endpoint

    // Scroll to bottom of chat when new messages are added
    useEffect(() => {
        if (chatContainerRef.current) {
            chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
        }
    }, [messages]);

    return (
        <div
            style={{
                display: 'flex',
                flexDirection: 'column',
                height: '100vh',
                backgroundColor: '#808080',
                fontFamily: 'Arial, sans-serif',
            }}
        >
            <header
                style={{
                    backgroundColor: '#000080',
                    color: '#FFFFFF',
                    padding: '10px',
                    textAlign: 'center',
                }}
            >
                <h1 style={{ fontSize: '24px', fontWeight: 'bold' }}>RAG Chat Demo</h1>
            </header>
            <main style={{ flex: '1', overflow: 'hidden', padding: '10px' }}>
                <div
                    ref={chatContainerRef}
                    style={{
                        height: '100%',
                        overflowY: 'auto',
                        border: '2px solid #000',
                        backgroundColor: '#FFFFFF',
                        padding: '10px',
                    }}
                >
                    {messages.map((msg, index) => (
                        <ChatMessage key={index} message={msg.text} isUser={msg.isUser} />
                    ))}
                    {isThinking && (
                        <div style={{ textAlign: 'center', color: '#000000' }}>
                            <p>Loading...</p>
                        </div>
                    )}
                </div>
            </main>
            <footer style={{ backgroundColor: '#C0C0C0', padding: '10px' }}>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                    <input
                        type="text"
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                        placeholder="Type your message..."
                        style={{
                            flex: '1',
                            padding: '5px',
                            marginRight: '5px',
                            border: '1px solid #000',
                            fontFamily: 'Arial, sans-serif',
                        }}
                    />
                    <button
                        onClick={sendMessage}
                        disabled={isThinking}
                        style={{
                            backgroundColor: '#000080',
                            color: '#FFFFFF',
                            padding: '5px 10px',
                            border: '1px solid #000',
                            cursor: 'pointer',
                            fontFamily: 'Arial, sans-serif',
                        }}
                    >
                        Send
                    </button>
                </div>
            </footer>
        </div>
    );
}
