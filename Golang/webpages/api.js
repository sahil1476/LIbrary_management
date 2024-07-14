const express = require('express');
const router = express.Router();
const User = require('../models/models.go/User'); // Assuming User model is in models/User.js

router.get('/', async (req, res) => {
  const { email, role } = req.query;

  try {
    const user = await User.findOne({ email });

    if (!user) {
      return res.status(404).json({ message: 'User not found' });
    }

    if (user.role === role) {
      if (role === 'admin') {
        res.redirect('/admin.html');
      } else if (role === 'user') {
        res.redirect('/user.html');
      } else {
        return res.status(403).json({ message: 'Invalid user role' });
      }
    } else {
      return res.status(403).json({ message: 'Invalid user role' });
    }
  } catch (error) {
    console.error(error);
    res.status(500).json({ message: 'Server error' });
  }
});

module.exports = router;