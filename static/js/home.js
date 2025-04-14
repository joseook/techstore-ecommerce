// Home page animations and interactions
document.addEventListener('DOMContentLoaded', function () {
  // Reveal animations on scroll
  const revealElements = document.querySelectorAll('.scroll-reveal');

  function checkReveal() {
    revealElements.forEach(element => {
      const elementTop = element.getBoundingClientRect().top;
      const windowHeight = window.innerHeight;

      if (elementTop < windowHeight - 100) {
        element.classList.add('visible');
      }
    });
  }

  // Initial check on page load
  checkReveal();

  // Check on scroll
  window.addEventListener('scroll', checkReveal);

  // Horizontal scroll controls
  const scrollLeftBtn = document.querySelector('.scroll-left');
  const scrollRightBtn = document.querySelector('.scroll-right');
  const scrollWrapper = document.querySelector('.scroll-wrapper');

  if (scrollLeftBtn && scrollRightBtn && scrollWrapper) {
    scrollLeftBtn.addEventListener('click', function () {
      scrollWrapper.scrollBy({
        left: -300,
        behavior: 'smooth'
      });
    });

    scrollRightBtn.addEventListener('click', function () {
      scrollWrapper.scrollBy({
        left: 300,
        behavior: 'smooth'
      });
    });
  }

  // Add to cart functionality
  const addToCartButtons = document.querySelectorAll('.add-to-cart');

  addToCartButtons.forEach(button => {
    button.addEventListener('click', function (e) {
      e.preventDefault();

      const productId = this.dataset.productId;
      const productName = this.dataset.productName;
      const productPrice = this.dataset.productPrice;

      // AJAX request to add to cart
      fetch('/cart/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRFToken': getCSRFToken()
        },
        body: JSON.stringify({
          product_id: productId,
          quantity: 1
        })
      })
        .then(response => response.json())
        .then(data => {
          if (data.success) {
            // Show success message
            showNotification(`${productName} added to cart`, 'success');

            // Update cart count
            updateCartCount(data.cart_count);
          } else {
            showNotification('Failed to add product to cart', 'error');
          }
        })
        .catch(error => {
          console.error('Error:', error);
          showNotification('An error occurred', 'error');
        });
    });
  });

  // Helper function to get CSRF token
  function getCSRFToken() {
    return document.querySelector('meta[name="csrf-token"]')?.getAttribute('content') || '';
  }

  // Helper function to show notification
  function showNotification(message, type) {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;

    document.body.appendChild(notification);

    // Fade in
    setTimeout(() => {
      notification.classList.add('show');
    }, 10);

    // Remove after 3 seconds
    setTimeout(() => {
      notification.classList.remove('show');

      // Remove element after fade out
      setTimeout(() => {
        notification.remove();
      }, 300);
    }, 3000);
  }

  // Helper function to update cart count
  function updateCartCount(count) {
    const cartCountElement = document.querySelector('.cart-count');
    if (cartCountElement) {
      cartCountElement.textContent = count;
      cartCountElement.classList.add('pulse');

      setTimeout(() => {
        cartCountElement.classList.remove('pulse');
      }, 500);
    }
  }
});

// Video background playback control
const heroVideo = document.querySelector('.hero-video');
if (heroVideo) {
  // Play the video when it's loaded
  heroVideo.addEventListener('loadeddata', function () {
    heroVideo.play().catch(error => {
      console.log('Auto-play was prevented:', error);
      // Add a play button as fallback if autoplay is blocked
      const playButton = document.createElement('button');
      playButton.className = 'video-play-button';
      playButton.innerHTML = '<i class="fas fa-play"></i>';
      playButton.addEventListener('click', () => heroVideo.play());
      heroVideo.parentNode.appendChild(playButton);
    });
  });

  // Pause video when not in viewport to save resources
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        heroVideo.play().catch(e => console.log('Play prevented:', e));
      } else {
        heroVideo.pause();
      }
    });
  }, { threshold: 0.2 });

  observer.observe(heroVideo);
} 