import React from "react";
import { Navigate } from "react-router-dom";

import { UserContext } from "../../contexts/userContext";

export default function Posts(props) {
    return (
        <div className='posts container space-y-2 flex flex-col justify-center items-center mx-auto w-auto'>    
            {props.children}
        </div>
    );
}

// export function Post({postID, displayName, avatar, content, date, topics, likes, isLiked}) {
export function Post({postData: {id, author: {display_name, avatar_url}, content, created_at, topics, likes, liked_by}, showLink}) {
    const {user} = React.useContext(UserContext);
    const [upvoteSelected, setUpvoteSelected] = React.useState(liked_by?.map(u => u.id).some(id => id === user.id));
    const [likesCount, setLikesCount] = React.useState(likes);
    const [disabled, setDisabled] = React.useState(false);

    if (created_at !== undefined && created_at !== null && created_at !== '') {
        created_at = new Date(created_at).toLocaleString();
    } else {
        created_at = 'Unknown Date';
    }

    React.useEffect(() => {
        const isLikedByUser = liked_by?.some(u => u.id === user.id);
        setUpvoteSelected(isLikedByUser);
    }, [liked_by, user.id]);

    const handlePostClick = () => {
        window.location.href = `/posts/${id}`;
    };

    const words = content.split(/(\s|\n)/);
    content = words.map((word, index) => {
        const urlRegex = /^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$/gm;
        if (urlRegex.test(word)) {
            return (
                <a className='text-blue-600 hover:underline hover:underline-offset-2' key={index} href={word.startsWith('http://') || word.startsWith('https://') ? word : 'https://' + word} target="_blank" rel="noopener noreferrer">{word}</a>
            );
        } else {
            return (
                <span key={index}>{word} </span>
            );
        }
    });

    return (
        <div className='post container rounded-md border-2 border-gray-300 w-1/2'>
            {
                showLink && 
                <div className='post-header text-xs flex flex-row text-wrap gap-3 p-2 bg-[#DDDDDD] w-full rounded-t hover:cursor-pointer' onClick={showLink ? handlePostClick : undefined}>
                    <img width={25} src={avatar_url} alt='avatar' />
                    <p>@{display_name ?? 'placeholder'}</p>
                    <p>{created_at}</p>
                </div>
            }
            {
                !showLink &&
                <div className='post-header text-xs flex flex-row text-wrap gap-3 p-2 bg-[#DDDDDD] w-full rounded-t'>
                    <img width={25} src={avatar_url} alt='avatar' />
                    <p>@{display_name ?? 'placeholder'}</p>
                    <p>{created_at}</p>
                </div>
            }
            {/* <hr /> */}
            <div className='post-content p-3 text-lg'>
                <p>{content ?? 'This is a placeholder post.'}</p>
                <ul className='flex flex-row gap-2 text-blue-500'>
                    {
                        topics?.map((topic, index) => (
                            <li key={index}>#{topic["name"]}</li>
                        ))
                    }
                </ul>
            </div>
            <div className='vote-footer p-3'>
                <div className='vote-footer-left flex flex-row gap-3'>
                    <button className={`vote-footer-left-like ${upvoteSelected ? 'bg-blue-200' : ''}`}
                        onClick={() => {
                            if (disabled) return;
                            if (upvoteSelected) {
                                fetch(`http://localhost:3030/posts/${id}/like/?type=unlike`, {
                                    method: 'POST', credentials: 'include'
                                }).then(response => response.json())
                                    .then(data => {
                                        setLikesCount(data.likes);
                                        setUpvoteSelected(false);
                                    })
                                    .catch(error => {
                                        console.log(error)
                                        setDisabled(true);
                                    });
                            } else {
                                fetch(`http://localhost:3030/posts/${id}/like/?type=like`, {
                                    method: 'POST', credentials: 'include'
                                }).then(response => response.json())
                                    .then(data => {
                                        setLikesCount(data.likes)
                                        setUpvoteSelected(true);
                                    })
                                    .catch(error => {
                                        console.log(error)
                                        setDisabled(true);
                                    });
                            }
                        }}
                    >
                        {/* <img width={30} height={30} src='/images/upvote.svg'/> */}
                        <p>👍</p>
                    </button>
                    <p>{likesCount ?? 0}</p>
                    {/* <p>💬</p> */}
                </div>
            </div>
        </div>
    );
}
