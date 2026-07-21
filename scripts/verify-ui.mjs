import { chromium } from 'playwright';

const URL = 'https://travel.sekaishopec.com';
const TIMEOUT = 15000;

(async () => {
  const browser = await chromium.launch({ headless: true });
  const ctx = await browser.newContext({
    viewport: { width: 390, height: 844 },
    deviceScaleFactor: 2,
  });
  const page = await ctx.newPage();
  
  console.log('[1] Navegando a', URL);
  await page.goto(URL, { waitUntil: 'networkidle', timeout: TIMEOUT });
  await page.waitForTimeout(3000);
  
  // Screenshot inicial
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/01-initial.png', fullPage: false });
  console.log('[1] Screenshot: initial state');

  // Verificar pin
  const pin = await page.$('#cytaxi-pin');
  console.log('[1] Pin visible:', !!pin);
  const pinBox = pin ? await pin.boundingBox() : null;
  console.log('[1] Pin boundingBox:', pinBox);

  // Verificar label
  const label = await page.$('#cytaxi-label');
  console.log('[1] Label visible:', !!label);
  const labelText = label ? await label.textContent() : null;
  console.log('[1] Label text:', labelText);

  // Verificar bottom sheet
  const sheet = await page.$('[aria-label="Panel de información del viaje"]');
  console.log('[1] Sheet visible:', !!sheet);
  const sheetBox = sheet ? await sheet.boundingBox() : null;
  console.log('[1] Sheet boundingBox:', sheetBox);
  
  // Verificar navbar
  const nav = await page.$('nav');
  console.log('[1] Navbar visible:', !!nav);
  const navBox = nav ? await nav.boundingBox() : null;
  console.log('[1] Navbar boundingBox:', navBox);

  // Verificar si hay espacio en blanco debajo del contenido
  const sheetContent = sheet ? await sheet.evaluate(el => ({
    scrollHeight: el.scrollHeight,
    clientHeight: el.clientHeight,
    offsetHeight: el.offsetHeight,
    children: el.children.length,
    innerDivHeight: el.querySelector('div:nth-child(2)')?.scrollHeight || 0,
  })) : null;
  console.log('[1] Sheet dimensions:', sheetContent);

  // Verificar si el pin está centrado en el mapa
  const mapContainer = await page.$('div[style*="position: absolute"]');
  const mapBox = mapContainer ? await mapContainer.boundingBox() : null;
  console.log('[1] Map container box:', mapBox);
  
  if (pinBox && mapBox) {
    const pinCenterX = pinBox.x + pinBox.width / 2;
    const mapCenterX = mapBox.x + mapBox.width / 2;
    console.log('[1] Pin X center:', pinCenterX, 'Map X center:', mapCenterX, 'Diff:', Math.abs(pinCenterX - mapCenterX));
  }

  // Buscar problemas de layout
  const bodyBg = await page.evaluate(() => {
    const body = document.body;
    return window.getComputedStyle(body).backgroundColor;
  });
  console.log('[1] Body background:', bodyBg);

  // Verificar si hay elementos overflow
  const overflowCheck = await page.evaluate(() => {
    const els = document.querySelectorAll('*');
    let overflows = [];
    for (const el of els) {
      const rect = el.getBoundingClientRect();
      if (rect.bottom > window.innerHeight + 10) {
        overflows.push({
          tag: el.tagName,
          class: el.className?.substring(0, 50),
          bottom: rect.bottom,
          windowH: window.innerHeight,
        });
      }
    }
    return overflows.slice(0, 5);
  });
  console.log('[1] Elements overflowing:', overflowCheck);

  // Verificar si el bottom sheet tiene espacio vacío
  const sheetWhitespace = await page.evaluate(() => {
    const sheet = document.querySelector('[aria-label="Panel de información del viaje"]');
    if (!sheet) return null;
    const rect = sheet.getBoundingClientRect();
    const content = sheet.querySelector('div:nth-child(2)');
    const contentRect = content?.getBoundingClientRect();
    return {
      sheetBottom: rect.bottom,
      sheetHeight: rect.height,
      contentBottom: contentRect?.bottom,
      contentHeight: contentRect?.height,
      whitespace: contentRect ? rect.bottom - contentRect.bottom : null,
    };
  });
  console.log('[1] Sheet whitespace:', sheetWhitespace);

  await browser.close();
  console.log('\n[DONE] Verificación completada. Revisa screenshots/');
})();
