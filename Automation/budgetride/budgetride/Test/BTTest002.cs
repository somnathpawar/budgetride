using System;
using System.Threading;
using budgetride.Pages;
using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace budgetride.Test
{
	[TestClass]
	public class BTTest002 : TestBase
    {
		[TestMethod]
        [TestProperty("Author", "Uttmesh Shukla")]
		[Description("Location can be entered sucessfully and cabs can be comapired & Booked ")]
        public void Enter_Pickup_and_DropPoint()
        {
			const string PickupPoint = "Brooklyn Bridge";
			const string DropPoint = "Wall Street";
			
			BaseApplicationPage baseApplicationPage = new BaseApplicationPage(Driver);
            HomePage homePage = new HomePage(Driver);

            baseApplicationPage.GoTO();

            homePage.LogoDisplayed();

			Driver.Manage().Window.FullScreen();

			homePage.EnterPickupAndDropPoint(PickupPoint, DropPoint);
                    
			homePage.ClickOnSearchCompareButton();

			Thread.Sleep(10000);

			homePage.Bookcab();
   
			Thread.Sleep(10000);

        }
    }
}   
