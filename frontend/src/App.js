import React from 'react';
import { Navigate } from 'react-router-dom';
import NavigationBar from './components/common/NavigationBar';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import HomePosts from './pages/HomePosts';
import HomeDiscussions from './pages/HomeDiscussions';
import Topics from './pages/Topics';

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <NavigationBar />
        <Routes>
          <Route path='/' element={<Navigate to="/posts" />} />
          <Route path='/posts' element={<HomePosts />} />
          <Route path='/discussions' element={<HomeDiscussions />} />
          <Route path='/topics' element={<Topics />} />
          <Route path='*' element={<h1>Not Found</h1>} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
