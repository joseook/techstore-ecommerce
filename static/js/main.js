document.addEventListener('DOMContentLoaded', () => {
  const currentPage = window.location.pathname;

  // Initialize animations
  initAnimations();

  if (currentPage === '/') {
    // Server-side rendering already loads the initial products
    // Just add the animation for revealing the sections
    initHeroAnimations();
    initPopularAndNewSections();
  } else if (currentPage === '/products') {
    loadAllProducts();
    loadFilters();
    setupSearchAndFilters();
  } else if (currentPage === '/support') {
    initSupportAnimations();
  }

  updateCartCount();
});

function initAnimations() {
  // Initialize scroll reveal animations
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
      }
    });
  }, {
    threshold: 0.1
  });

  document.querySelectorAll('.scroll-reveal').forEach(el => observer.observe(el));
}

function initHeroAnimations() {
  const heroContent = document.querySelector('.hero-content');
  if (heroContent) {
    setTimeout(() => {
      heroContent.classList.add('visible');
    }, 100);
  }
}

function initSupportAnimations() {
  const supportCards = document.querySelectorAll('.support-card');
  const faqItems = document.querySelectorAll('.faq-item');

  supportCards.forEach((card, index) => {
    setTimeout(() => {
      card.classList.add('visible');
    }, index * 200);
  });

  faqItems.forEach((item, index) => {
    setTimeout(() => {
      item.classList.add('visible');
    }, index * 200);
  });
}

// Initializes the animations for the popular and new product sections
function initPopularAndNewSections() {
  const popularSection = document.querySelector('.popular-products');
  const newSection = document.querySelector('.new-products');

  if (popularSection) {
    setTimeout(() => {
      popularSection.classList.add('visible');
    }, 300);
  }

  if (newSection) {
    setTimeout(() => {
      newSection.classList.add('visible');
    }, 500);
  }
}

async function loadPopularProducts() {
  try {
    const response = await fetch('/api/products?popular=true');
    const products = await response.json();
    const container = document.getElementById('popular-products');

    products.forEach((product, index) => {
      const productCard = createProductCard(product);
      container.appendChild(productCard);

      setTimeout(() => {
        productCard.classList.add('visible');
      }, index * 100);
    });
  } catch (error) {
    console.error('Error loading popular products:', error);
  }
}

async function loadNewProducts() {
  try {
    const response = await fetch('/api/products?new=true');
    const products = await response.json();
    const container = document.getElementById('new-products');

    products.forEach((product, index) => {
      const productCard = createProductCard(product);
      container.appendChild(productCard);

      setTimeout(() => {
        productCard.classList.add('visible');
      }, index * 100);
    });
  } catch (error) {
    console.error('Error loading new products:', error);
  }
}

/**
 * Carrega todos os produtos da API
 */
async function loadAllProducts() {
  const productsContainer = document.getElementById('productsContainer');
  const loadingIndicator = document.getElementById('loadingIndicator');

  // Limpar container e mostrar indicador de carregamento
  if (productsContainer) {
    console.log("Container de produtos encontrado");
    productsContainer.innerHTML = '';

    if (loadingIndicator) {
      loadingIndicator.style.display = 'flex';
      console.log("Indicador de carregamento ativado");
    }

    try {
      console.log("Iniciando carregamento de produtos...");
      const response = await fetch('/api/products');

      if (!response.ok) {
        throw new Error(`Erro ao carregar produtos: ${response.status} ${response.statusText}`);
      }

      const products = await response.json();
      console.log(`Produtos carregados: ${products.length} itens`);

      if (products.length === 0) {
        productsContainer.innerHTML = '<div class="col-12 text-center"><p>Nenhum produto encontrado.</p></div>';
      } else {
        // Criar uma row para os produtos
        const row = document.createElement('div');
        row.className = 'row';
        productsContainer.appendChild(row);

        // Renderizar cada produto
        console.log("Renderizando produtos...");
        products.forEach(product => {
          const productCol = document.createElement('div');
          productCol.className = 'col-lg-3 col-md-4 col-sm-6 product-item';
          productCol.innerHTML = createProductCard(product);
          row.appendChild(productCol);
        });

        // Adicionar event listeners aos botões de ação
        setupProductActions();

        console.log("Produtos renderizados com sucesso!");
      }
    } catch (error) {
      console.error("Erro ao carregar produtos:", error);
      productsContainer.innerHTML = `<div class="col-12 text-center"><p>Erro ao carregar produtos: ${error.message}</p></div>`;
    } finally {
      if (loadingIndicator) {
        loadingIndicator.style.display = 'none';
        console.log("Indicador de carregamento desativado");
      }
    }
  } else {
    console.error("Container de produtos não encontrado");
  }
}

/**
 * Cria o HTML do card de produto
 * @param {Object} product - Dados do produto
 * @returns {string} HTML do card de produto
 */
function createProductCard(product) {
  // Normalizar dados do produto
  const name = product.name || 'Produto sem nome';
  const price = product.price || 0;
  const oldPrice = product.old_price || null;
  const imageUrl = product.image || '/static/img/product-placeholder.jpg';
  const inStock = product.in_stock === undefined ? true : product.in_stock;
  const isNew = product.is_new === true;
  const isDiscount = oldPrice !== null && oldPrice > price;

  return `
    <div class="product-card h-100 visible" style="opacity: 1; transform: translateY(0);">
      <div class="product-image">
        <img src="${imageUrl}" alt="${name}" class="img-fluid">
        <div class="product-actions">
          <button class="action-button quick-view" data-id="${product.id || 0}">
            <i class="bi bi-eye"></i>
          </button>
          <button class="action-button add-to-wishlist" data-id="${product.id || 0}">
            <i class="bi bi-heart"></i>
          </button>
          <button class="action-button add-to-cart" onclick="addToCart('${product.id || 0}')">
            <i class="bi bi-cart-plus"></i>
          </button>
        </div>
        ${isNew ? '<span class="badge bg-success position-absolute top-0 start-0 m-3">New</span>' : ''}
        ${isDiscount ? '<span class="badge bg-danger position-absolute top-0 end-0 m-3">-' + Math.round((1 - price / oldPrice) * 100) + '%</span>' : ''}
      </div>
      <div class="product-info p-3">
        <h3 class="product-title">
          <a href="/product/${product.id || 0}">${name}</a>
        </h3>
        <div class="d-flex justify-content-between align-items-center">
          <span class="product-price">$${price.toFixed(2)}</span>
          <div class="product-rating">
            <i class="bi bi-star-fill"></i>
            <i class="bi bi-star-fill"></i>
            <i class="bi bi-star-fill"></i>
            <i class="bi bi-star-fill"></i>
            <i class="bi bi-star-half"></i>
            <span class="rating-count">(24)</span>
          </div>
        </div>
        <p class="product-description d-none">${product.description || ''}</p>
      </div>
    </div>
  `;
}

/**
 * Configura os event listeners para os botões de ação dos produtos
 */
function setupProductActions() {
  // Quick view buttons
  document.querySelectorAll('.quick-view').forEach(button => {
    button.addEventListener('click', function () {
      const productId = this.getAttribute('data-id');

      // Get product details via API
      fetch(`/api/products/${productId}`)
        .then(response => response.json())
        .then(product => {
          // Populate modal with product details
          const modal = document.getElementById('quickViewModal');
          if (!modal) return;

          document.getElementById('quickViewTitle').textContent = 'Quick View: ' + product.name;
          document.getElementById('quickViewName').textContent = product.name;
          document.getElementById('quickViewPrice').textContent = `$${product.price.toFixed(2)}`;
          document.getElementById('quickViewDescription').textContent = product.description;
          document.getElementById('quickViewCategory').textContent = product.category;
          document.getElementById('quickViewImage').src = product.image_url || product.imageURL;

          // Set product ID for "Add to Cart" button
          document.getElementById('modalAddToCart').setAttribute('data-id', product.id);

          // Reset quantity input
          document.getElementById('quantityInput').value = 1;

          // Show the modal
          const quickViewModal = new bootstrap.Modal(modal);
          quickViewModal.show();
        })
        .catch(error => {
          console.error('Error loading product details:', error);
          showToast('Failed to load product details', 'error');
        });
    });
  });

  // Add to cart buttons
  document.querySelectorAll('.add-to-cart').forEach(button => {
    button.addEventListener('click', function (e) {
      e.preventDefault();
      e.stopPropagation();

      const productId = this.getAttribute('data-id') || this.getAttribute('onclick')?.match(/'([^']+)'/)?.[1];
      if (!productId) return;

      // Animate button
      this.classList.add('adding');

      // Add product to cart
      addToCart(productId);

      // Show success toast
      showToast('Product added to cart!', 'success');

      // Reset button state after animation
      setTimeout(() => {
        this.classList.remove('adding');
      }, 1000);
    });
  });

  // Add to wishlist buttons
  document.querySelectorAll('.add-to-wishlist').forEach(button => {
    button.addEventListener('click', function (e) {
      e.preventDefault();
      e.stopPropagation();

      const productId = this.getAttribute('data-id');
      if (!productId) return;

      // Toggle wishlist state
      this.classList.toggle('active');

      if (this.classList.contains('active')) {
        // Add to wishlist
        this.innerHTML = '<i class="bi bi-heart-fill"></i>';
        showToast('Added to wishlist!', 'success');
      } else {
        // Remove from wishlist
        this.innerHTML = '<i class="bi bi-heart"></i>';
        showToast('Removed from wishlist', 'info');
      }
    });
  });
}

/**
 * Exibe uma notificação toast
 * @param {string} message - Mensagem a ser exibida
 * @param {string} type - Tipo de toast (success, error, info)
 */
function showToast(message, type = 'info') {
  const toast = document.createElement('div');
  toast.className = `toast ${type}`;
  toast.textContent = message;
  document.body.appendChild(toast);

  // Remover o toast após a animação
  setTimeout(() => {
    document.body.removeChild(toast);
  }, 3500);
}

/**
 * Adiciona um produto ao carrinho
 * @param {string} productId - ID do produto
 */
function addToCart(productId, quantity = 1) {
  console.log(`Adicionando produto ID ${productId} ao carrinho, quantidade: ${quantity}`);

  fetch('/cart/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      product: { id: productId },
      quantity: quantity
    })
  })
    .then(response => response.json())
    .then(data => {
      // Mostra notificação de sucesso
      showToast('Product added to cart!', 'success');

      // Atualiza o contador do carrinho
      updateCartCount();

      // Abre o carrinho lateral
      updateCartSidebar();
    })
    .catch(error => {
      console.error('Error adding to cart:', error);
      showToast('Failed to add product to cart', 'error');
    });
}

/**
 * Remove um produto do carrinho
 * @param {string} productId - ID do produto
 */
function removeFromCart(productId) {
  fetch('/cart/remove', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      product: { id: productId }
    })
  })
    .then(response => response.json())
    .then(data => {
      showToast('Product removed from cart', 'info');
      updateCartCount();
      updateCartSidebar();
    })
    .catch(error => {
      console.error('Error removing from cart:', error);
      showToast('Failed to remove product from cart', 'error');
    });
}

/**
 * Atualiza o contador de itens no carrinho
 */
function updateCartCount() {
  fetch('/cart')
    .then(response => response.json())
    .then(data => {
      const count = data.cart.items.reduce((total, item) => total + item.quantity, 0);
      const cartCountElements = document.querySelectorAll('.cart-count');

      cartCountElements.forEach(element => {
        element.textContent = count;
        element.style.display = count > 0 ? 'inline-flex' : 'none';
      });
    })
    .catch(error => {
      console.error('Error updating cart count:', error);
    });
}

/**
 * Atualiza o conteúdo do carrinho lateral
 * @param {boolean} open - Se deve abrir o carrinho lateral após atualizar
 */
function updateCartSidebar(open = true) {
  // Verificar se o elemento do carrinho lateral existe
  const cartSidebar = document.getElementById('cartSidebar');
  if (!cartSidebar) return;

  fetch('/cart')
    .then(response => response.json())
    .then(data => {
      const container = document.getElementById('cart-items-container');
      const emptyMessage = document.getElementById('empty-cart-message');
      const cartSummary = document.getElementById('cart-summary');

      // Limpar itens anteriores
      if (container.children.length > 0 && container.children[0].id !== 'empty-cart-message') {
        container.innerHTML = '';
      }

      if (data.cart.items.length === 0) {
        // Carrinho vazio
        if (!document.getElementById('empty-cart-message')) {
          container.innerHTML = `
            <div class="text-center py-5" id="empty-cart-message">
              <i class="bi bi-cart text-muted" style="font-size: 3rem;"></i>
              <p class="mt-3">Your cart is empty</p>
              <a href="/products" class="btn btn-primary mt-3">Continue Shopping</a>
            </div>
          `;
        }
        if (cartSummary) cartSummary.style.display = 'none';
      } else {
        // Carrinho com itens
        container.innerHTML = ''; // Limpar conteúdo

        let subtotal = 0;

        // Adicionar cada item ao carrinho
        data.cart.items.forEach(item => {
          const itemTotal = item.product.price * item.quantity;
          subtotal += itemTotal;

          const itemElement = document.createElement('div');
          itemElement.className = 'cart-item d-flex align-items-center mb-3 p-2 border-bottom';
          itemElement.innerHTML = `
            <img src="${item.product.image || item.product.imageURL}" alt="${item.product.name}" class="img-fluid me-3" style="width: 60px; height: 60px; object-fit: cover; border-radius: 4px;">
            <div class="flex-grow-1">
              <h6 class="mb-0">${item.product.name}</h6>
              <div class="d-flex justify-content-between align-items-center">
                <span class="text-muted">${item.quantity} × $${item.product.price.toFixed(2)}</span>
                <span>$${itemTotal.toFixed(2)}</span>
              </div>
            </div>
            <button class="btn btn-sm text-danger" onclick="removeFromCart('${item.product.id}')">
              <i class="bi bi-x-circle"></i>
            </button>
          `;
          container.appendChild(itemElement);
        });

        // Atualizar o sumário
        if (cartSummary) {
          cartSummary.style.display = 'block';
          const tax = subtotal * 0.1; // 10% de imposto
          const total = subtotal + tax;

          document.getElementById('cart-subtotal').textContent = `$${subtotal.toFixed(2)}`;
          document.getElementById('cart-tax').textContent = `$${tax.toFixed(2)}`;
          document.getElementById('cart-total').textContent = `$${total.toFixed(2)}`;
        }

        // Abrir o carrinho lateral se solicitado
        if (open) {
          const bsOffcanvas = new bootstrap.Offcanvas(cartSidebar);
          bsOffcanvas.show();
        }
      }
    })
    .catch(error => {
      console.error('Error updating cart sidebar:', error);
    });
}

/**
 * Adiciona um produto à lista de desejos
 * @param {string} productId - ID do produto
 */
function addToWishlist(productId) {
  console.log(`Adicionando produto ID ${productId} à lista de desejos`);

  // Simulação de adição à lista de desejos (em uma aplicação real, isso seria uma requisição ao servidor)
  const button = document.querySelector(`.add-to-wishlist[data-id="${productId}"]`);
  if (button) {
    button.classList.toggle('active');

    if (button.classList.contains('active')) {
      button.innerHTML = '<i class="bi bi-heart-fill"></i>';
      showToast('Added to wishlist!', 'success');
    } else {
      button.innerHTML = '<i class="bi bi-heart"></i>';
      showToast('Removed from wishlist', 'info');
    }
  } else {
    showToast('Added to wishlist!', 'success');
  }
}

/**
 * Compartilha um produto
 * @param {string} productId - ID do produto
 * @param {string} productName - Nome do produto
 */
function shareProduct(productId, productName) {
  const url = window.location.origin + '/product/' + productId;
  const title = productName || 'Check out this product';

  if (navigator.share) {
    navigator.share({
      title: title,
      url: url
    })
      .then(() => console.log('Successful share'))
      .catch((error) => {
        console.log('Error sharing:', error);
        fallbackShare(url);
      });
  } else {
    fallbackShare(url);
  }
}

/**
 * Método alternativo de compartilhamento quando a Web Share API não está disponível
 * @param {string} url - URL para compartilhar
 */
function fallbackShare(url) {
  // Copia para a área de transferência
  navigator.clipboard.writeText(url)
    .then(() => showToast('Link copied to clipboard!', 'success'))
    .catch(err => {
      console.error('Error copying text: ', err);

      // Criar um campo de texto temporário
      const textarea = document.createElement('textarea');
      textarea.value = url;
      textarea.style.position = 'fixed';
      document.body.appendChild(textarea);
      textarea.focus();
      textarea.select();

      try {
        document.execCommand('copy');
        showToast('Link copied to clipboard!', 'success');
      } catch (err) {
        showToast('Failed to copy link', 'error');
      }

      document.body.removeChild(textarea);
    });
}

// Inicializar o contador do carrinho quando o DOM estiver pronto
document.addEventListener('DOMContentLoaded', function () {
  updateCartCount();

  // Configurar o botão do carrinho na navbar para abrir o carrinho lateral
  const cartButtons = document.querySelectorAll('.cart-icon, .nav-link[href="/cart"]');
  cartButtons.forEach(button => {
    button.addEventListener('click', function (e) {
      // Se o elemento do carrinho lateral existir, previnir o comportamento padrão
      // e abrir o carrinho lateral em vez de navegar para /cart
      if (document.getElementById('cartSidebar')) {
        e.preventDefault();
        updateCartSidebar(true);
      }
    });

    // Atualizar atributos para funcionar com o Bootstrap Offcanvas
    if (document.getElementById('cartSidebar')) {
      button.setAttribute('data-bs-toggle', 'offcanvas');
      button.setAttribute('data-bs-target', '#cartSidebar');
      button.setAttribute('aria-controls', 'cartSidebar');
      button.removeAttribute('href');
    }
  });
});

async function loadFilters() {
  try {
    // Check if we're on a page that needs filters
    const categoryFilter = document.getElementById('categoryFilter');
    const brandFilter = document.getElementById('brandFilter');

    if (!categoryFilter || !brandFilter) {
      // Skip if we're not on the products page
      return;
    }

    // Load categories
    const categoriesResponse = await fetch('/api/categories');
    const categories = await categoriesResponse.json();

    categories.forEach(category => {
      const option = document.createElement('option');
      option.value = category;
      option.textContent = category;
      categoryFilter.appendChild(option);
    });

    // Load brands
    const brandsResponse = await fetch('/api/brands');
    const brands = await brandsResponse.json();

    brands.forEach(brand => {
      const option = document.createElement('option');
      option.value = brand;
      option.textContent = brand;
      brandFilter.appendChild(option);
    });
  } catch (error) {
    console.error('Error loading filters:', error);
  }
}

function setupSearchAndFilters() {
  const searchInput = document.getElementById('searchInput');
  const searchButton = document.getElementById('searchButton');
  const categoryFilter = document.getElementById('categoryFilter');
  const brandFilter = document.getElementById('brandFilter');
  const sortOption = document.getElementById('sortOption');

  // Skip if we're not on the products page
  if (!searchInput) {
    return;
  }

  let timeout;
  // Set products per page to 8
  const PRODUCTS_PER_PAGE = 8;
  // Initialize with page 1
  let currentPage = 1;
  // Store filtered products
  let filteredProducts = [];

  async function updateProducts() {
    clearTimeout(timeout);

    // Show loading indicator
    const loadingIndicator = document.getElementById('loadingIndicator');
    if (loadingIndicator) {
      loadingIndicator.style.display = 'flex';
    }

    timeout = setTimeout(async () => {
      const search = searchInput?.value || '';
      const category = categoryFilter?.value || '';
      const brand = brandFilter?.value || '';
      const sort = sortOption?.value || 'popular';

      const params = new URLSearchParams();
      if (search) params.append('search', search);
      if (category) params.append('category', category);
      if (brand) params.append('brand', brand);
      params.append('sort', sort);

      try {
        const response = await fetch(`/api/products?${params.toString()}`);
        const products = await response.json();
        const container = document.getElementById('productsContainer');

        if (!container) {
          console.log('Products container not found for filtering');
          return; // Skip if container doesn't exist
        }

        // Store filtered products for pagination
        filteredProducts = products;

        // Reset to first page when filters change
        currentPage = 1;

        // Render the current page
        renderCurrentPage();

      } catch (error) {
        console.error('Error updating products:', error);
        if (loadingIndicator) {
          loadingIndicator.style.display = 'none';
        }
      }
    }, 300);
  }

  function renderCurrentPage() {
    const container = document.getElementById('productsContainer');
    const loadingIndicator = document.getElementById('loadingIndicator');

    if (!container) return;

    // Clear container but keep loading indicator
    while (container.firstChild) {
      if (container.firstChild.id === 'loadingIndicator') {
        break;
      }
      container.removeChild(container.firstChild);
    }

    if (filteredProducts.length === 0) {
      // Create row for "no products found" message
      const row = document.createElement('div');
      row.className = 'row';
      container.appendChild(row);

      // Add "no products found" message
      const noProductsCol = document.createElement('div');
      noProductsCol.className = 'col-12 text-center py-5';
      noProductsCol.innerHTML = `
        <div class="alert alert-info">
          <h4>No products found</h4>
          <p>Try adjusting your filters or search criteria.</p>
        </div>
      `;
      row.appendChild(noProductsCol);
    } else {
      // Create row for products
      const row = document.createElement('div');
      row.className = 'row g-4';
      container.appendChild(row);

      // Calculate start and end index for current page
      const startIndex = (currentPage - 1) * PRODUCTS_PER_PAGE;
      const endIndex = Math.min(startIndex + PRODUCTS_PER_PAGE, filteredProducts.length);

      // Display only products for current page
      const currentPageProducts = filteredProducts.slice(startIndex, endIndex);

      // Add products to row
      currentPageProducts.forEach(product => {
        const productCol = document.createElement('div');
        productCol.className = 'col-md-6 col-lg-4 col-xl-3 product-item';
        productCol.dataset.category = product.category;
        productCol.dataset.price = product.price;
        productCol.dataset.id = product.id;

        productCol.innerHTML = `
          <div class="product-card h-100">
            <div class="product-image">
              <img src="${product.image_url || product.imageURL}" alt="${product.name}" class="img-fluid">
              <div class="product-actions">
                <button class="action-button quick-view" data-id="${product.id}">
                  <i class="bi bi-eye"></i>
                </button>
                <button class="action-button add-to-wishlist" data-id="${product.id}">
                  <i class="bi bi-heart"></i>
                </button>
                <button class="action-button add-to-cart" onclick="addToCart('${product.id}')">
                  <i class="bi bi-cart-plus"></i>
                </button>
              </div>
              ${product.is_new ? '<span class="badge bg-success position-absolute top-0 start-0 m-3">New</span>' : ''}
              ${product.is_popular ? '<span class="badge bg-primary position-absolute top-0 end-0 m-3">Popular</span>' : ''}
            </div>
            <div class="product-info p-3">
              <h3 class="product-title">
                <a href="/product/${product.id}">${product.name}</a>
              </h3>
              <div class="d-flex justify-content-between align-items-center">
                <span class="product-price">$${product.price.toFixed(2)}</span>
                <div class="product-rating">
                  ${getRatingStars(product.rating)}
                  <span class="rating-count">(${Math.floor(Math.random() * 50) + 5})</span>
                </div>
              </div>
              <p class="product-description d-none">${product.description}</p>
            </div>
          </div>
        `;

        row.appendChild(productCol);
      });

      // Update pagination
      updatePagination();
    }

    // Hide loading indicator
    if (loadingIndicator) {
      loadingIndicator.style.display = 'none';
    }

    // Setup action buttons for filtered products
    setupProductActions();

    // Setup quick view for filtered products
    setupQuickView();
  }

  // Helper function to generate star rating HTML
  function getRatingStars(rating) {
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating % 1 >= 0.5;
    const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0);

    let starsHTML = '';

    // Add full stars
    for (let i = 0; i < fullStars; i++) {
      starsHTML += '<i class="bi bi-star-fill"></i>';
    }

    // Add half star if needed
    if (hasHalfStar) {
      starsHTML += '<i class="bi bi-star-half"></i>';
    }

    // Add empty stars
    for (let i = 0; i < emptyStars; i++) {
      starsHTML += '<i class="bi bi-star"></i>';
    }

    return starsHTML;
  }

  // Update pagination based on filtered products
  function updatePagination() {
    const paginationElement = document.querySelector('.pagination');
    if (!paginationElement) return;

    const pageCount = Math.ceil(filteredProducts.length / PRODUCTS_PER_PAGE);

    paginationElement.innerHTML = '';

    // Previous button
    const prevLi = document.createElement('li');
    prevLi.className = 'page-item' + (currentPage === 1 ? ' disabled' : '');
    prevLi.innerHTML = `<a class="page-link" href="#" tabindex="-1">Previous</a>`;
    prevLi.addEventListener('click', function (e) {
      e.preventDefault();
      if (currentPage > 1) {
        currentPage--;
        renderCurrentPage();
      }
    });
    paginationElement.appendChild(prevLi);

    // Page numbers
    for (let i = 1; i <= pageCount; i++) {
      const pageLi = document.createElement('li');
      pageLi.className = 'page-item' + (i === currentPage ? ' active' : '');
      pageLi.innerHTML = `<a class="page-link" href="#">${i}</a>`;
      pageLi.addEventListener('click', function (e) {
        e.preventDefault();
        currentPage = i;
        renderCurrentPage();
      });
      paginationElement.appendChild(pageLi);
    }

    // Next button
    const nextLi = document.createElement('li');
    nextLi.className = 'page-item' + (currentPage >= pageCount ? ' disabled' : '');
    nextLi.innerHTML = `<a class="page-link" href="#">Next</a>`;
    nextLi.addEventListener('click', function (e) {
      e.preventDefault();
      if (currentPage < pageCount) {
        currentPage++;
        renderCurrentPage();
      }
    });
    paginationElement.appendChild(nextLi);
  }

  // Add event listeners
  if (searchInput) {
    searchInput.addEventListener('input', updateProducts);
  }

  if (searchButton) {
    searchButton.addEventListener('click', updateProducts);
  }

  if (categoryFilter) {
    categoryFilter.addEventListener('change', updateProducts);
  }

  if (brandFilter) {
    brandFilter.addEventListener('change', updateProducts);
  }

  if (sortOption) {
    sortOption.addEventListener('change', updateProducts);
  }

  // Load products on init
  updateProducts();
}

// Update order summary
function updateOrderSummary(cart) {
  if (!cart.items) return;

  const subtotal = cart.items.reduce((sum, item) => sum + (item.Product.Price * item.Quantity), 0);
  const shipping = subtotal > 0 ? 5.99 : 0;
  const total = subtotal + shipping;

  document.getElementById('subtotal').textContent = `$${subtotal.toFixed(2)}`;
  document.getElementById('shipping').textContent = `$${shipping.toFixed(2)}`;
  document.getElementById('total').textContent = `$${total.toFixed(2)}`;
}

// Update quantity
function updateQuantity(productId, change) {
  fetch('/cart/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      product: { id: productId },
      quantity: change
    })
  })
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to update cart');
      }
      return response.json();
    })
    .then(data => {
      updateCartCount();
      showToast('Cart updated!');
    })
    .catch(error => {
      console.error('Error updating quantity:', error);
      showToast('Error updating cart');
    });
}

// Proceed to checkout
function proceedToCheckout() {
  window.location.href = '/checkout';
}

// Setup quick view functionality for dynamically created products
function setupQuickView() {
  const quickViewButtons = document.querySelectorAll('.quick-view');
  if (quickViewButtons.length === 0) return;

  const quickViewModal = new bootstrap.Modal(document.getElementById('quickViewModal'));

  quickViewButtons.forEach(button => {
    button.addEventListener('click', function () {
      const productId = this.getAttribute('data-id');
      const productItem = document.querySelector(`.product-item[data-id="${productId}"]`);

      if (!productItem) return;

      const productCard = productItem.querySelector('.product-card');
      const productName = productCard.querySelector('.product-title a').innerText;
      const productPrice = productCard.querySelector('.product-price').innerText;
      const productImage = productCard.querySelector('.product-image img').src;
      const productDescription = productCard.querySelector('.product-description').innerText;
      const productCategory = productItem.getAttribute('data-category');

      // Set modal content
      document.getElementById('quickViewTitle').innerText = productName;
      document.getElementById('quickViewName').innerText = productName;
      document.getElementById('quickViewPrice').innerText = productPrice;
      document.getElementById('quickViewImage').src = productImage;
      document.getElementById('quickViewDescription').innerText = productDescription;
      document.getElementById('quickViewCategory').innerText = productCategory;

      // Set add to cart button functionality
      document.getElementById('modalAddToCart').setAttribute('data-id', productId);

      quickViewModal.show();
    });
  });

  // Set up quantity input handlers in quick view modal
  const decreaseBtn = document.getElementById('decreaseQuantity');
  const increaseBtn = document.getElementById('increaseQuantity');
  const quantityInput = document.getElementById('quantityInput');

  if (decreaseBtn && increaseBtn && quantityInput) {
    decreaseBtn.addEventListener('click', function () {
      let value = parseInt(quantityInput.value);
      if (value > 1) {
        quantityInput.value = value - 1;
      }
    });

    increaseBtn.addEventListener('click', function () {
      let value = parseInt(quantityInput.value);
      quantityInput.value = value + 1;
    });

    // Add to cart from modal
    const modalAddToCartBtn = document.getElementById('modalAddToCart');
    if (modalAddToCartBtn) {
      modalAddToCartBtn.addEventListener('click', function () {
        const productId = this.getAttribute('data-id');
        const quantity = parseInt(quantityInput.value);

        addToCart(productId, quantity);
        quickViewModal.hide();
      });
    }
  }
}

function updateProducts(products) {
  const container = document.getElementById('productsContainer');
  if (!container) return;

  // Clear container except loading indicator
  const loadingIndicator = document.getElementById('loadingIndicator');
  container.innerHTML = '';

  // Re-append loading indicator if it existed
  if (loadingIndicator) {
    loadingIndicator.style.display = 'none';
    container.appendChild(loadingIndicator);
  }

  if (products.length === 0) {
    // Display no products found message
    const noProductsMsg = document.createElement('div');
    noProductsMsg.className = 'col-12 text-center py-5';
    noProductsMsg.innerHTML = '<h3>No products found</h3><p>Try changing your filters or search criteria.</p>';
    container.appendChild(noProductsMsg);
    return;
  }

  products.forEach((product, index) => {
    const productCard = createProductCard(product);
    container.appendChild(productCard);

    // Make products visible immediately
    productCard.classList.add('visible');
  });

  // Setup quick view for dynamically created products
  setupQuickView();
} 