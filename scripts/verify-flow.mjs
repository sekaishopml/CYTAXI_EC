import { chromium } from 'playwright';

const URL = 'https://travel.sekaishopec.com';
const TIMEOUT = 15000;

(async () => {
  const browser = await chromium.launch({ headless: true });
  const ctx = await browser.newContext({
    viewport: { width: 390, height: 844 },
    deviceScaleFactor: 2,
    permissions: ['geolocation'],
    geolocation: { latitude: -2.1894, longitude: -79.8893 },
  });
  const page = await ctx.newPage();

  console.log('[1] === ESTADO INICIAL: pickup_select ===');
  await page.goto(URL, { waitUntil: 'networkidle', timeout: TIMEOUT });
  await page.waitForTimeout(3500);
  
  // Medir whitespace
  const sheetInfo1 = await page.evaluate(() => {
    const sheet = document.querySelector('[aria-label="Panel de información del viaje"]');
    if (!sheet) return null;
    const sheetRect = sheet.getBoundingClientRect();
    const contentDiv = sheet.querySelector('div:nth-child(2)');
    const contentRect = contentDiv?.getBoundingClientRect();
    return {
      sheetH: Math.round(sheetRect.height),
      contentH: Math.round(contentRect?.height || 0),
      whitespace: Math.round(sheetRect.bottom - (contentRect?.bottom || 0)),
    };
  });
  console.log('[1] Sheet whitespace:', sheetInfo1);

  // Medir pin position vs center
  const pinInfo1 = await page.evaluate(() => {
    const wrapper = document.querySelector('#cytaxi-pin')?.parentElement;
    const wrapperRect = wrapper?.getBoundingClientRect();
    return {
      pinX: Math.round((wrapperRect?.x || 0) + (wrapperRect?.width || 0) / 2),
      pinY: Math.round((wrapperRect?.y || 0) + (wrapperRect?.height || 0)),
      windowCenter: Math.round(window.innerWidth / 2),
    };
  });
  console.log('[1] Pin centering:', pinInfo1);
  console.log('[1] Pin centered?', Math.abs(pinInfo1.pinX - pinInfo1.windowCenter) < 5 ? 'YES ✓' : 'NO ✗');

  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/01-pickup-select.png' });

  // Click "Confirmar ubicación"
  console.log('\n[2] === CLICK CONFIRMAR UBICACIÓN → input state ===');
  const confirmBtn = await page.getByRole('button', { name: /Confirmar ubicación/i });
  if (confirmBtn) {
    await confirmBtn.click();
    await page.waitForTimeout(1000);
  }
  
  const state2 = await page.evaluate(() => {
    const label = document.querySelector('#cytaxi-label');
    return {
      labelText: label?.textContent || null,
      labelVisible: !!label,
    };
  });
  console.log('[2] Label:', state2);

  const sheetInfo2 = await page.evaluate(() => {
    const sheet = document.querySelector('[aria-label="Panel de información del viaje"]');
    if (!sheet) return null;
    const sheetRect = sheet.getBoundingClientRect();
    const contentDiv = sheet.querySelector('div:nth-child(2)');
    const contentRect = contentDiv?.getBoundingClientRect();
    return {
      sheetH: Math.round(sheetRect.height),
      contentH: Math.round(contentRect?.height || 0),
      whitespace: Math.round(sheetRect.bottom - (contentRect?.bottom || 0)),
    };
  });
  console.log('[2] Sheet whitespace:', sheetInfo2);
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/02-input-state.png' });

  // Verificar pin B azul
  const pinB = await page.evaluate(() => {
    const paths = document.querySelectorAll('#cytaxi-pin svg path');
    return paths.length > 0 ? paths[0].getAttribute('fill') : null;
  });
  console.log('[2] Pin fill (should be blue gradient):', pinB);

  // Buscar destino
  console.log('\n[3] === BUSCAR DESTINO ===');
  const destInput = await page.getByRole('textbox', { name: /Buscar destino/i });
  if (destInput) {
    await destInput.fill('Malabo, Ecuador');
    await page.waitForTimeout(1500);
  }
  
  // Verificar sugerencias
  const suggestions = await page.evaluate(() => {
    const items = document.querySelectorAll('button[aria-label^="Seleccionar"]');
    return items.length;
  });
  console.log('[3] Suggestions found:', suggestions);
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/03-search-dest.png' });

  // Click primera sugerencia
  if (suggestions > 0) {
    const firstSug = await page.getByRole('button', { name: /Seleccionar/i }).first();
    await firstSug?.click();
    await page.waitForTimeout(500);
    console.log('[3] Selected destination');
    await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/04-dest-selected.png' });
  }

  // Click "Buscar viaje"
  console.log('\n[4] === CONFIRMAR DESTINO → confirm state ===');
  const searchBtn = await page.getByRole('button', { name: /Buscar viaje/i });
  if (searchBtn) {
    await searchBtn.click();
    await page.waitForTimeout(4000); // Esperar route calculation
  }

  const sheetInfo4 = await page.evaluate(() => {
    const sheet = document.querySelector('[aria-label="Panel de información del viaje"]');
    if (!sheet) return null;
    const sheetRect = sheet.getBoundingClientRect();
    const contentDiv = sheet.querySelector('div:nth-child(2)');
    const contentRect = contentDiv?.getBoundingClientRect();
    return {
      sheetH: Math.round(sheetRect.height),
      contentH: Math.round(contentRect?.height || 0),
      whitespace: Math.round(sheetRect.bottom - (contentRect?.bottom || 0)),
    };
  });
  console.log('[4] Sheet whitespace (confirm):', sheetInfo4);
  await page.screenshot({ path: '/home/CYTAXI_EC/screenshots/05-confirm-state.png' });

  // Verificar que hay ruta dibujada en el mapa
  const hasRoute = await page.evaluate(() => {
    const svgs = document.querySelectorAll('.leaflet-pane path, .leaflet-overlay-pane path');
    return svgs.length;
  });
  console.log('[4] SVG path elements on map:', hasRoute);

  console.log('\n[DONE] Todos los screenshots en /home/CYTAXI_EC/screenshots/');
  await browser.close();
})();
