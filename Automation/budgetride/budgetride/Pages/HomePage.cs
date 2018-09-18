using System;
using System.Threading;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using OpenQA.Selenium;

namespace budgetride.Pages
{
	public class HomePage : BaseApplicationPage
    {
        public HomePage(IWebDriver driver) : base(driver)
        {
        }

        /*---- Web elements of Home page ----*/
		private IWebElement AppLogo => Driver.FindElement(By.XPath("//img[@class='app_logo']"));
		private IWebElement PickupTextbox => Driver.FindElement(By.Id("origin-input"));
		private IWebElement DroptoTextbox => Driver.FindElement(By.Id("destination-input"));
		private IWebElement SearchCompareButton => Driver.FindElement(By.Id("search_compare"));
		private IWebElement ResetButton => Driver.FindElement(By.Id("reset_form"));
		private IWebElement Bookbuttonone => Driver.FindElement(By.XPath("//tbody//tr[2]//td[4]//button[1]//span[1]"));
		private IWebElement Bookbuttontwo => Driver.FindElement(By.XPath("//tbody//tr[3]//td[4]//button[1]"));
		private IWebElement Bookbuttonthree => Driver.FindElement(By.XPath("//tbody//tr[3]//td[4]//button[1]"));
		private IWebElement FullScreen => Driver.FindElement(By.XPath("//button[@title='Toggle fullscreen view']"));



		// Check for Logo displayed method 
        internal void LogoDisplayed()
        {
            Assert.IsTrue(AppLogo.Displayed);

        }

		//Search & Compare button click method 
		internal void ClickOnSearchCompareButton()
        {
			Assert.IsTrue(SearchCompareButton.Displayed);
			SearchCompareButton.Click();
        }


		//Reset button click method 
		internal void ClickOnResetButton(String ButtonName)
        {
            
			Assert.IsTrue(ResetButton.Displayed);
			ResetButton.Click();
        }

		      
		// Enter values into Pickup & Drop point text boxes method 
		internal void EnterPickupAndDropPoint (string Pickuppoint, string Droppoint)
		{
   
			PickupTextbox.Click();
			PickupTextbox.Clear();
			PickupTextbox.SendKeys(Pickuppoint);
			Thread.Sleep(2000);
			PickupTextbox.SendKeys(Keys.ArrowDown);
			PickupTextbox.SendKeys(Keys.Enter);

			PickupTextbox.Click();
           	DroptoTextbox.SendKeys(Droppoint);
			Thread.Sleep(2000);
			DroptoTextbox.SendKeys(Keys.ArrowDown);
			DroptoTextbox.SendKeys(Keys.Enter);
			Thread.Sleep(3000);
   
		}

        // Click on book cab button method 
		internal void Bookcab()
        {
			Assert.IsTrue(Bookbuttonone.Displayed);
			Bookbuttonone.Click();

        }
    }
}
